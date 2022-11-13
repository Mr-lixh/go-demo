package huaweicloud_demo

import (
	"fmt"
	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/region"
	"testing"
)

var (
	imsClient *ims.ImsClient

	jobId = "ff8080817f8e255e0180adc1d52e034f"
)

func init() {
	imsClient = ims.NewImsClient(ims.ImsClientBuilder().WithRegion(region.ValueOf(regionName)).WithCredential(auth).Build())
}

/*
* ShowJob: 查询Job状态

*** 镜像管理：
* ListImages: 查询镜像列表
* CreateImage: 制作镜像，可以从ecs实例转化
*
*** 标签管理：
* AddImageTag: 添加镜像标签
* BatchAddOrDeleteTags: 批量添加删除镜像标签
* DeleteImageTag: 删除镜像标签
* ListImageByTags: 按标签查询镜像
 */

func TestShowImageJob(t *testing.T) {
	request := &model.ShowJobRequest{
		JobId: jobId,
	}

	response, err := imsClient.ShowJob(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}

// 查询镜像列表
func TestListImages(t *testing.T) {
	imagetypeRequest := model.GetListImagesRequestImagetypeEnum().PRIVATE
	isregisteredRequest := model.GetListImagesRequestIsregisteredEnum().TRUE
	tagRequest := "app=iat"

	request := &model.ListImagesRequest{
		Imagetype:    &imagetypeRequest,
		Isregistered: &isregisteredRequest,
		Tag:          &tagRequest,
	}

	response, err := imsClient.ListImages(request)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", response.HttpStatusCode)
	for _, i := range *response.Images {
		// image id: fe72a014-e0eb-40fb-b1a2-aeee0e925a00, image name: iat-test
		fmt.Printf("image id: %s, image name: %s, image tags: %+v\n", i.Id, i.Name, i.Tags)
	}
}

// 从ecs实例创建镜像
func TestCreateImage(t *testing.T) {
	request := &model.CreateImageRequest{}

	var listImageTagsbody = []model.TagKeyValue{
		{
			Key:   "app",
			Value: "iat",
		},
		{
			Key:   "comp",
			Value: "neartv",
		},
		{
			Key:   "deployid",
			Value: "20220510123456",
		},
		{
			Key:   "version",
			Value: "latest",
		},
	}
	instanceIdCreateImageRequestBody := serverId
	descriptionCreateImageRequestBody := "iat image"
	request.Body = &model.CreateImageRequestBody{
		Name:        "iat-test",
		InstanceId:  &instanceIdCreateImageRequestBody,
		ImageTags:   &listImageTagsbody,
		Description: &descriptionCreateImageRequestBody,
	}

	response, err := imsClient.CreateImage(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}

func TestListImageByTags(t *testing.T) {
	request := &model.ListImageByTagsRequest{}

	var listValuesTags = []model.Tags{
		{
			Key:    "app",
			Values: []string{"iat"},
		},
		{
			Key:    "comp",
			Values: []string{"neartv"},
		},
		{
			Key:    "version",
			Values: []string{"latest"},
		},
	}

	request.Body = &model.ListImageByTagsRequestBody{
		Tags:   &listValuesTags,
		Action: model.GetListImageByTagsRequestBodyActionEnum().FILTER,
	}

	response, err := imsClient.ListImageByTags(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}
