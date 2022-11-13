package huaweicloud_demo

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/region"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"testing"
)

/*
* ListSubnets: 获取子网列表
 */

var (
	vpcClient *vpc.VpcClient
)

func init() {
	vpcClient = vpc.NewVpcClient(vpc.VpcClientBuilder().WithRegion(region.ValueOf(regionName)).WithCredential(auth).Build())
}

func TestListVpcs(t *testing.T) {
	request := &model.ListVpcsRequest{}

	response, err := vpcClient.ListVpcs(request)
	if err != nil {
		t.Fatal(err)
	}

	for _, vpc := range *response.Vpcs {
		// vpc id: a67bf70f-d851-4b0c-9bab-b793a9d729db, vpc name: vpc-test
		fmt.Printf("vpc id: %s, vpc name: %s", vpc.Id, vpc.Name)
	}
}

func TestListSubnets(t *testing.T) {
	request := &model.ListSubnetsRequest{
		VpcId: &vpcId,
	}

	response, err := vpcClient.ListSubnets(request)
	if err != nil {
		t.Fatal(err)
	}

	for _, sn := range *response.Subnets {
		// subnet id: ba647982-f46e-4909-bbcb-60c0d9202a32, subnet name: subnet-5abc
		fmt.Printf("subnet id: %s, subnet name: %s", sn.Id, sn.Name)
	}
}
