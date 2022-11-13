package controllers

import (
	"context"

	"fmt"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

const (
	PodCostFinalizer        = "aicloud.iflytek.com/pod-cost"
	MaxConcurrentReconciles = 10
)

type PodCostReconciler struct {
	client.Client

	Namespaces string
}

func (r *PodCostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	pod := &corev1.Pod{}
	err := r.Get(ctx, req.NamespacedName, pod)

	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("pod resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get pod")
		return ctrl.Result{}, err
	}

	//if pod.ObjectMeta.DeletionTimestamp.IsZero() {
	//	if err := r.insertFinalizerIfMissing(ctx, logger, pod, PodCostFinalizer); err != nil {
	//		return ctrl.Result{}, err
	//	}
	//} else {
	//	// TODO: 如果 pod 被删除，需要处理具体的业务逻辑--pod运行数据回写到运营数据库
	//	if err := r.removeFinalizerIfExists(ctx, logger, pod, PodCostFinalizer); err != nil {
	//		return ctrl.Result{}, err
	//	}
	//}

	return ctrl.Result{}, nil
}

func (r *PodCostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("pod_cost").
		For(&corev1.Pod{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(object client.Object) bool {
			if r.Namespaces == "" {
				return false
			}

			// TODO: 根据用户配置，只对部分命名空间下的pod标注finalizer
			if object.GetNamespace() == r.Namespaces {
				return true
			}

			return false
		})).
		Complete(r)
}

func (r *PodCostReconciler) removeFinalizerIfExists(ctx context.Context, log logr.Logger, instance *corev1.Pod, finalizerName string) error {
	if !controllerutil.ContainsFinalizer(instance, finalizerName) {
		return nil
	}

	log.Info("Write pod cost to db")
	if err := r.writeToDB(ctx, instance); err != nil {
		log.Error(err, "Failed to write to db")
		return err
	}
	log.Info("Writed pod cost to db")

	log.Info("Removing pod-cost finalizer")
	clone := instance.DeepCopy()
	clone.SetFinalizers(removeString(clone.GetFinalizers(), finalizerName))
	if err := r.Update(ctx, clone); err != nil {
		log.Error(err, "Failed to remove pod-cost finalizer")
		return err
	}
	log.Info("Removed pod-cost finalizer")
	return nil
}

func (r *PodCostReconciler) insertFinalizerIfMissing(ctx context.Context, log logr.Logger, instance *corev1.Pod, finalizerName string) error {
	if controllerutil.ContainsFinalizer(instance, finalizerName) {
		return nil
	}

	log.Info("Inserting pod-cost finalizer")
	clone := instance.DeepCopy()
	clone.SetFinalizers(append(clone.GetFinalizers(), finalizerName))
	if err := r.Update(ctx, clone); err != nil {
		log.Error(err, "Failed to insert pod-cost finalizer")
		return err
	}
	log.Info("Inserted pod-cost finalizer")

	return nil
}

func (r *PodCostReconciler) writeToDB(ctx context.Context, instance *corev1.Pod) error {
	// TODO: 将 pod 运行数据写回运营数据库
	// TODO: 如果运行时长小于一定阈值，直接返回，防止 pod 频繁驱逐造成的多次计量
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("pod run %d seconds", time.Now().Unix()-instance.CreationTimestamp.Unix()))
	logger.Info(fmt.Sprintf("pod resource: %d cpu, %d memory", instance.Spec.Containers[0].Resources.Requests.Cpu().Value(), instance.Spec.Containers[0].Resources.Requests.Memory().Value()))

	// FIXME: 建议这里添加重试，如果数据库写入失败，重试 N 次后，直接返回 nil,
	// 重试失败后，上报数据到 metrics.
	// 如果使用 workqueue 的重试机制，会一直重试，导致 pod 一直 terminating.
	ErrPodResource.WithLabelValues(instance.Namespace, instance.Name, corev1.ResourceCPU.String()).Set(instance.Spec.Containers[0].Resources.Requests.Cpu().AsApproximateFloat64())

	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}

	return
}
