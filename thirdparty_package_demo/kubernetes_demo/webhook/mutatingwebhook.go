package main

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type HybridCloud struct {
	Client  client.Client
	decoder *admission.Decoder
}

// HybridCloud only print the pod
func (h *HybridCloud) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger := log.FromContext(ctx)

	pod := &corev1.Pod{}

	err := h.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// TODO: Do some business logic.
	logger.Info("monitoring pod", "pod", pod.Name)

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// HybridCloud implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (h *HybridCloud) InjectDecoder(d *admission.Decoder) error {
	h.decoder = d
	return nil
}
