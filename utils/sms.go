package utils

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/inhzus/go-berater/config"
	"log"
)

func SendSMS(phone, code string) error {
	c := config.GetConfig()
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", c.Code.AccessKey, c.Code.AccessSecret)
	if err != nil {
		return err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = "南大咨询"
	request.QueryParams["TemplateCode"] = "SMS_163433313"
	request.QueryParams["TemplateParam"] = "{\"code\":\"" + code + "\"}"

	res, err := client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}
	log.Print(res)
	var j map[string]interface{}
	_ = json.Unmarshal(res.GetHttpContentBytes(), &j)
	if j["Code"] == "OK" {
		return nil
	} else {
		return errors.New(j["Message"].(string))
	}
}
