package spy

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/golang/glog"
)

func ConfigHTTPClient(config *Config) {
	resty.SetAllowGetMethodPayload(true)
	resty.SetQueryParams(config.GlobalSettings.Params)
	resty.SetHeaders(config.GlobalSettings.Headers)
	resty.SetAuthToken(config.GlobalSettings.Authtoken)
}

func DoTest(test TestCase, host string) {
	// Create request
	request := resty.R()
	request.SetQueryParams(test.Params)
	request.SetMultiValueQueryParams(test.MultiParams)
	request.SetHeaders(test.Headers)
	request.SetBody(test.Body)
	request.SetAuthToken(test.Authtoken)
	request.SetFormData(test.Form)
	request.SetMultiValueFormData(test.MultiForm)
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
		fmt.Printf("\nError: %v", err)
	} else {
		glog.Infof("\nResponse Body: %v\nDuration: %v", response, response.Time())
	}
	glog.Flush()
}

func Dotests(testCases []TestCase, host string) {
	for _, test := range testCases {
		DoTest(test, host)
	}
}
