package main

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	iai "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/iai/v20180301"
)

//https://console.cloud.tencent.com/api/explorer?Product=iai&Version=2018-03-01&Action=DetectFace

func main() {

	credential := common.NewCredential(
		"AKIDGP5g6LjI0EcprmxmZGfdP3WzRkruFdHl",
		"3TinlZw9B2lNsCKLBzgmFZif5JLBdv2x",
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "iai.tencentcloudapi.com"
	client, _ := iai.NewClient(credential, "ap-beijing", cpf)

	request := iai.NewDetectFaceRequest()

	params := `{"Url":"https://wx1.sinaimg.cn/mw690/6ae35d94gy1g36z7198q4j21w02io1ky.jpg"}`
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.DetectFace(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", response.ToJsonString())
}
