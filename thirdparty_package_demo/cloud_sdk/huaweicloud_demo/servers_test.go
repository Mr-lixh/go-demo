package huaweicloud_demo

import (
	"fmt"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/region"
	"testing"
)

var (
	ecsClient *ecs.EcsClient

	ecsJobId = "ff8080817f8e255e0180ae77893c2fc4"
)

func init() {
	ecsClient = ecs.NewEcsClient(ecs.EcsClientBuilder().WithRegion(region.ValueOf(regionName)).WithHttpConfig(httpConfig).WithCredential(auth).Build())
}

/*
*** 生命周期管理
* CreatePostPaidServers: 创建云服务器(按需)
* ListServersDetails: 查询云服务器详情列表
* DeleteServers: 删除云服务器
*
*** 批量操作
* BatchAttachSharableVolumes: 批量挂载指定共享卷
*
*** 可用区管理
* NovaListAvailabilityZones: 查询可用区列表
 */

func TestShowEcsJob(t *testing.T) {
	request := &model.ShowJobRequest{
		JobId: ecsJobId,
	}

	response, err := ecsClient.ShowJob(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}

func TestCreatePostPaidServers(t *testing.T) {
	request := &model.CreatePostPaidServersRequest{}

	var listServerTagsServer = []model.PostPaidServerTag{
		{
			Key:   "app",
			Value: "iat",
		},
	}
	rootVolumeServer := &model.PostPaidServerRootVolume{
		Volumetype: model.GetPostPaidServerRootVolumeVolumetypeEnum().SAS,
	}
	var listNicsServer = []model.PostPaidServerNic{
		{
			SubnetId: subnetId,
		},
	}
	// 设置管理用户密码，使用密码登陆方式不支持实例自定义数据注入功能。
	//adminPassServerPostPaidServer:= "lixiaohe@@52"
	// 注入用户数据。参考规范：https://support.huaweicloud.com/usermanual-ecs/zh-cn_topic_0032380449.html
	// 查看自定义数据：curl http://169.254.169.254/openstack/latest/user_data
	userDataServerPostPaidServer := "IyEvYmluL2Jhc2gKZXhwb3J0IGNsdXN0ZXJfbmFtZT1keCAmJiBiYXNoIC9yb290L3NjcmlwdHMvaW5pdF9ub2RlLnNoIHwgdGVlIC90bXAvaW5pdF9ub2RlLmxvZwo="
	countServer := int32(1)
	serverBody := &model.PostPaidServer{
		FlavorRef:  flavorRef,
		ImageRef:   imageRef,
		Name:       "iat-test",
		Count:      &countServer,
		Nics:       listNicsServer,
		RootVolume: rootVolumeServer,
		ServerTags: &listServerTagsServer,
		//AdminPass: &adminPassServerPostPaidServer,
		UserData: &userDataServerPostPaidServer,
		Vpcid:    vpcId,
	}
	dryRunCreatePostPaidServersRequestBody := true
	request.Body = &model.CreatePostPaidServersRequestBody{
		Server: serverBody,
		DryRun: &dryRunCreatePostPaidServersRequestBody,
	}

	response, err := ecsClient.CreatePostPaidServers(request)
	if err != nil {
		t.Fatal(err)
	}

	// CreatePostPaidServersResponse {"job_id":"ff8080817f8e188b0180a9a50bf86a00","serverIds":["82034d06-c0b4-43cf-a1a9-e9295d9c2e6b"]}
	fmt.Printf("%+v\n", response)
}

func TestListServersDetails(t *testing.T) {
	//request := &model.ListServersDetailsRequest{}
	//
	//tagRequest := "app=iat"
	//request.Tags = &tagRequest
	//
	//response, err := ecsClient.ListServersDetails(request)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Printf("%+v\n", response)
	request := &model.NovaListServersDetailsRequest{}
	response, err := ecsClient.NovaListServersDetails(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}

func TestDeleteServers(t *testing.T) {
	request := &model.DeleteServersRequest{}

	var listServersBody = []model.ServerId{
		{
			Id: serverId,
		},
	}
	deleteVolumeDeleteServersRequestBody := true
	deletePublicipDeleteServersRequestBody := true

	request.Body = &model.DeleteServersRequestBody{
		DeleteVolume:   &deleteVolumeDeleteServersRequestBody,
		DeletePublicip: &deletePublicipDeleteServersRequestBody,
		Servers:        listServersBody,
	}

	response, err := ecsClient.DeleteServers(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}

func TestBatchAttachSharableVolumes(t *testing.T) {
	request := &model.BatchAttachSharableVolumesRequest{
		VolumeId: volumeId,
	}

	deviceServerinfoBatchAttachSharableVolumesOption := "/data1"
	var listServerinfobody = []model.BatchAttachSharableVolumesOption{
		{
			ServerId: serverId,
			Device:   &deviceServerinfoBatchAttachSharableVolumesOption,
		},
	}

	request.Body = &model.BatchAttachSharableVolumesRequestBody{
		Serverinfo: listServerinfobody,
	}

	response, err := ecsClient.BatchAttachSharableVolumes(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}

func TestNovaListAvailabilityZones(t *testing.T) {
	request := &model.NovaListAvailabilityZonesRequest{}

	response, err := ecsClient.NovaListAvailabilityZones(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}
