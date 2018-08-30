package spy

import (
	"github.com/go-resty/resty"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"time"
)

func ConfigHTTPClient(client *resty.Client, APIsetting *TestCase, Clientsetting *ClientSetting) {
	client.SetAllowGetMethodPayload(true)
	client.SetQueryParams(APIsetting.Params)
	client.SetHeaders(APIsetting.Headers)

	if APIsetting.AuthToken != "" {
		client.SetAuthToken(APIsetting.AuthToken)
	}
	if APIsetting.BasicAuth.Username != "" {
		client.SetBasicAuth(APIsetting.BasicAuth.Username, APIsetting.BasicAuth.Password)
	}

	if Clientsetting.RetryCount > 0 {
		client.SetRetryCount(Clientsetting.RetryCount)
		if Clientsetting.RetryWait>0{
			client.SetRetryWaitTime(time.Duration(Clientsetting.RetryWait) * time.Millisecond)
		}
		if Clientsetting.RetryMaxWait>0{
			client.SetRetryMaxWaitTime(time.Duration(Clientsetting.RetryMaxWait) * time.Millisecond)
		}
	}

	if Clientsetting.Timeout != 0 {
		client.SetTimeout(time.Duration(Clientsetting.Timeout) * time.Millisecond)
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

	glog.V(1).Infof("host %s, method %s, url %s", host, test.Method, test.URL)
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
		glog.V(1).Infof("Request fail, Duration: %v", response.Time())
	} else {
		glog.V(1).Infof("Request success:\n%s\n Duration: %v", response, response.Time())
	}
	glog.Flush()
	return err, response
}

func Dotests(clientset *kubernetes.Clientset, config *Config, service *VictimService, chaos *Chaos) {

	for _, testcases := range config.TestCaseLists {
		client := resty.New()
		// Apply global http client settings
		ConfigHTTPClient(client, &config.APISetting, &config.ClientSetting)
		// Apply local http client settings
		ConfigHTTPClient(client, &testcases.APIsetting, &testcases.ClientSetting)
		// Find host
		var host string
		if testcases.Host == "" {
			service, err := GetService(clientset, config, testcases.Service)
			if err != nil {
				glog.Errorf("Can't find service \"%s\", continue: %s", testcases.Service, err)
				continue
			}
			host = service.Spec.ClusterIP
		} else {
			host = testcases.Host
		}
		// Do tests
		for _, test := range testcases.TestCases {
			err, response := DoTest(client, test, host)
			AddResponse(service, chaos, &test, response, err)
		}
		// Send response to db
		SendResponses()
	}

}
