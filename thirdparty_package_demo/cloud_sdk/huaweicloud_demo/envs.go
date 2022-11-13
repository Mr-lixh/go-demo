package huaweicloud_demo

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	"net/http"
)

var (
	ak         = "EGIFYQDPOZINAQYWOVNW"
	sk         = "iLUOyAxxhs24nMXQkJbYGUjUnqE8iP02Z3qQBMyY"
	regionName = "cn-east-2"
	//projectId = "423c17ae1b3e4c5a8eb748f9c2bb64d6"
	vpcId     = "a67bf70f-d851-4b0c-9bab-b793a9d729db"
	subnetId  = "ba647982-f46e-4909-bbcb-60c0d9202a32"
	flavorRef = "s6.small.1"
	imageRef  = "1b87702a-e46a-4217-b7db-f70d6ea8322f" // iat-test Build from CentOS 7.8 64bit
	serverId  = "b0817f97-0995-4d51-bc09-b74b9dce0a7a"
	volumeId  = "" // 共享磁盘ID
)

var (
	auth       basic.Credentials
	httpConfig *config.HttpConfig
)

func init() {
	auth = basic.NewCredentialsBuilder().WithAk(ak).WithSk(sk).Build()
	httpConfig = config.DefaultHttpConfig().WithIgnoreSSLVerification(true).
		WithHttpHandler(httphandler.NewHttpHandler().
			AddRequestHandler(RequestHandler).
			AddResponseHandler(ResponseHandler))
}

func RequestHandler(request http.Request) {
	fmt.Println(request)
}

func ResponseHandler(response http.Response) {
	fmt.Println(response)
}
