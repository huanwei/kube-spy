package spy

import (
	"github.com/golang/glog"
	client_v2 "github.com/influxdata/influxdb/client/v2"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var DBClient client_v2.Client
var respBP client_v2.BatchPoints
var pingBP client_v2.BatchPoints

// Connnect to Database
func ConnectDB(clientset *kubernetes.Clientset,config *Config) {
	// Get Database address
	pods,err:=clientset.CoreV1().Pods(config.Namespace).List(meta_v1.ListOptions{LabelSelector:"name=influxdb-spy"})
	if err!=nil{
		glog.Fatalf("Fail to find database pod: %s",err)
		panic(err)
	}

	DBClient, err = client_v2.NewHTTPClient(client_v2.HTTPConfig{
		Addr:     "http://"+pods.Items[0].Status.PodIP+":8086",
		Username: "kubespy",
		Password: "kubespy",
	})
	if err != nil {
		glog.Fatalf("Fail to connect database: %s", err)
		panic(err)
	}

	// Create response points batch
	respBP, err = client_v2.NewBatchPoints(client_v2.BatchPointsConfig{
		Database:  "spy",
		Precision: "ms",
	})
	if err != nil {
		glog.Fatalf("Fail to create points batch: %s", err)
		panic(err)
	}

}

func AddResponse(url, method, body, duration string) {
	// Create map
	tags := make(map[string]string)
	fileds := make(map[string]interface{})

	// Set tags and fields
	tags["url"]=url
	tags["method"]=method
	fileds["body"] = body
	fileds["duration"]=duration

	// Create point
	point, err := client_v2.NewPoint(
		"response",
		tags,
		fileds,
	)
	if err != nil {
		glog.Warningf("Fail to create point: %s", err)
	} else {
		// Add to batch
		respBP.AddPoint(point)
	}
}

func SendResponses() {
	var err error
	// Write batch
	DBClient.Write(respBP)
	// Create new batch
	respBP, err = client_v2.NewBatchPoints(client_v2.BatchPointsConfig{
		Database:  "spy",
		Precision: "ms",
	})
	if err != nil {
		glog.Fatalf("Fail to create points batch: %s", err)
		panic(err)
	}
}

func AddPingResult(url,method,delay,loss string){
	// Create map
	tags := make(map[string]string)
	fileds := make(map[string]interface{})

	// TODO:Set tags and fields here
	tags["url"]=url
	tags["method"]=method
	fileds["delay"] = delay
	fileds["loss"]=loss

	// Create point
	point, err := client_v2.NewPoint(
		"ping",
		tags,
		fileds,
	)
	if err != nil {
		glog.Warningf("Fail to create point: %s", err)
	} else {
		// Add to batch
		pingBP.AddPoint(point)
	}
}

func SendPingResults() {
	var err error
	// Write batch
	DBClient.Write(pingBP)
	// Create new batch
	pingBP, err = client_v2.NewBatchPoints(client_v2.BatchPointsConfig{
		Database:  "spy",
		Precision: "ms",
	})
	if err != nil {
		glog.Fatalf("Fail to create points batch: %s", err)
		panic(err)
	}
}