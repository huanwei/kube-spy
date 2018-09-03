package pkg

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/golang/glog"
	"time"
)

func SendRequest(config RequestConfig, host string) (response *resty.Response) {
	resty.SetTimeout(2000*time.Millisecond)
	request := resty.R()
	request.SetQueryParams(config.Params)
	request.SetMultiValueQueryParams(config.MultiParams)
	request.SetHeaders(config.Headers)
	request.SetBody(config.Body)
	request.SetAuthToken(config.Authtoken)
	request.SetFormData(config.Form)
	request.SetMultiValueFormData(config.MultiForm)
	if config.PathParams != nil {
		request.SetPathParams(config.PathParams)
	}
	if config.Files != nil {
		request.SetFiles(config.Files)
	}

	glog.Infof("method %s", config.Method)
	var err error
	// Select method and send
	switch config.Method {
	case "GET", "Get", "get":
		{
			response, err = request.Get("http://" + host + config.URL)
		}
	case "POST", "Post", "post":
		{
			response, err = request.Post("http://" + host + config.URL)
		}
	case "PUT", "Put", "put":
		{
			response, err = request.Put("http://" + host + config.URL)
		}
	case "PATCH", "Patch", "patch":
		{
			response, err = request.Patch("http://" + host + config.URL)
		}
	case "DELETE", "Delete", "delete":
		{
			response, err = request.Delete("http://" + host + config.URL)
		}
	case "OPTIONS", "Options", "options":
		{
			response, err = request.Options("http://" + host + config.URL)
		}
	case "HEAD", "Head", "head":
		{
			response, err = request.Head("http://" + host + config.URL)
		}
	default:
		glog.Warningf("Unsupported http method: %s, skip this test case", config.Method)
	}

	// Check potential error
	if err != nil {
		fmt.Printf("\nError: %v", err)
	} else {
		glog.Infof("\nResponse Body: %v\nDuration: %v", response, response.Time())
	}
	glog.Flush()
	if err != nil {
		glog.Errorf("http get err:%s", err)
	}

	return response
}
