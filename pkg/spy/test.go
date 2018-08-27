package spy

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/golang/glog"
	"time"
)

func ConfigHTTPClient(client *resty.Client, config *Config) {
	client.SetAllowGetMethodPayload(true)
	client.SetQueryParams(config.APISetting.Params)
	client.SetHeaders(config.APISetting.Headers)

	if config.APISetting.AuthToken != "" {
		client.SetAuthToken(config.APISetting.AuthToken)
	}
	if config.APISetting.BasicAuth.Username != "" {
		client.SetBasicAuth(config.APISetting.BasicAuth.Username, config.APISetting.BasicAuth.Password)
	}

	if config.ClientSetting.RetryCount > 0 {
		client.SetRetryCount(config.ClientSetting.RetryCount)
		client.SetRetryWaitTime(time.Duration(config.ClientSetting.RetryWait) * time.Millisecond)
		client.SetRetryMaxWaitTime(time.Duration(config.ClientSetting.RetryMaxWait) * time.Millisecond)
	}

	if config.ClientSetting.Timeout != 0 {
		client.SetTimeout(time.Duration(config.ClientSetting.Timeout) * time.Millisecond)
	}

}

func DoTest(client *resty.Client, test TestCase, host string) (error, *resty.Response) {
	// Create request
	request := client.R()
	request.SetQueryParams(test.Params)
	request.SetMultiValueQueryParams(test.MultiParams)
	request.SetHeaders(test.Headers)
	request.SetFormData(test.Form)
	request.SetMultiValueFormData(test.MultiForm)

	if test.Body != "" {
		request.SetBody(test.Body)
	}
	if test.AuthToken != "" {
		request.SetAuthToken(test.AuthToken)
	}
	if test.BasicAuth.Username != "" {
		request.SetBasicAuth(test.BasicAuth.Username, test.BasicAuth.Password)
	}
	if test.PathParams != nil {
		request.SetPathParams(test.PathParams)
	}
	if test.Files != nil {
		request.SetFiles(test.Files)
	}

	var response *resty.Response

	glog.Infof("method %s", test.Method)
	var err error
	// Select method and send
	switch test.Method {
	case "GET", "Get", "get":
		{
			response, err = request.Get("http://" + host + test.URL)
		}
	case "POST", "Post", "post":
		{
			response, err = request.Post("http://" + host + test.URL)
		}
	case "PUT", "Put", "put":
		{
			response, err = request.Put("http://" + host + test.URL)
		}
	case "PATCH", "Patch", "patch":
		{
			response, err = request.Patch("http://" + host + test.URL)
		}
	case "DELETE", "Delete", "delete":
		{
			response, err = request.Delete("http://" + host + test.URL)
		}
	case "OPTIONS", "Options", "options":
		{
			response, err = request.Options("http://" + host + test.URL)
		}
	case "HEAD", "Head", "head":
		{
			response, err = request.Head("http://" + host + test.URL)
		}
	default:
		glog.Warningf("Unsupported http method: %s, skip this test case", test.Method)
	}

	// Check potential error
	if err != nil {
		glog.Infof("Request fail, Duration: %v", response.Time())
	} else {
		glog.Infof("Request success:\n%s\n Duration: %v", response, response.Time())
	}
	glog.Flush()
	return err, response
}

//func NewRestyClient(config *Config)*resty.Client{
//	client := resty.New()
//	ConfigHTTPClient(client, config)
//	return client
//}

func Dotests(config *Config, host string, service *VictimService, chaos *Chaos) {
	client := resty.New()
	ConfigHTTPClient(client, config)
	for _, test := range config.TestCases {
		err, response := DoTest(client, test, host)
		AddResponse(service, chaos, &test, response, err)
	}
	// Send response to db
	SendResponses()
}
