package spy

import (
	"encoding/base64"
	"github.com/go-resty/resty"
	"github.com/golang/glog"
	client_v2 "github.com/influxdata/influxdb/client/v2"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

var DBClient client_v2.Client
var respBP client_v2.BatchPoints
var pingBP client_v2.BatchPoints

// Connnect to Database
func ConnectDB(clientset *kubernetes.Clientset, config *Config) {
	// Get Database address
	pods, err := clientset.CoreV1().Pods(config.Namespace).List(meta_v1.ListOptions{LabelSelector: "name=influxdb-spy"})
	if err != nil {
		glog.Fatalf("Fail to find database pod: %s", err)
		glog.Flush()
		panic(err)
	}

	DBClient, err = client_v2.NewHTTPClient(client_v2.HTTPConfig{
		Addr:     "http://" + pods.Items[0].Status.PodIP + ":8086",
		Username: "kubespy",
		Password: "kubespy",
	})
	if err != nil {
		glog.Fatalf("Fail to connect database: %s", err)
		glog.Flush()
		panic(err)
	}

	// Create response points batch
	respBP, err = client_v2.NewBatchPoints(client_v2.BatchPointsConfig{
		Database:  "spy",
		Precision: "ms",
	})
	if err != nil {
		glog.Fatalf("Fail to create points batch: %s", err)
		glog.Flush()
		panic(err)
	}

	// Create ping points batch
	pingBP, err = client_v2.NewBatchPoints(client_v2.BatchPointsConfig{
		Database:  "spy",
		Precision: "ms",
	})
	if err != nil {
		glog.Fatalf("Fail to create points batch: %s", err)
		glog.Flush()
		panic(err)
	}
}

func AddResponse(service *VictimService, chaos *Chaos, test *TestCase, response *resty.Response, err error) {
	// Create map
	tags := make(map[string]string)
	fields := make(map[string]interface{})

	// Set tags and fields
	if service == nil {
		tags["victim"] = "none"
	} else {
		tags["victim"] = service.Name
	}
	tags["url"] = test.URL
	tags["method"] = test.Method

	if chaos == nil {
		fields["chaos-ingress"] = "none"
		fields["chaos-egress"] = "none"
		fields["chaos-replica"] = "none"
	} else {
		fields["chaos-ingress"] = chaos.Ingress
		fields["chaos-egress"] = chaos.Egress
		if chaos.Replica == 0 {
			fields["chaos-replica"] = "none"
		} else {
			fields["chaos-replica"] = strconv.Itoa(chaos.Replica)
		}
	}

	fields["status"] = response.Status()

	if err != nil {
		fields["body"] = base64.StdEncoding.EncodeToString([]byte(err.Error()))
	} else {
		fields["body"] = base64.StdEncoding.EncodeToString([]byte(response.Body()))
	}

	fields["duration"] = response.Time()

	// Create point
	point, err := client_v2.NewPoint(
		"response",
		tags,
		fields,
	)
	if err != nil {
		glog.Warningf("Fail to create point: %s", err)
	} else {
		// Add to batch
		respBP.AddPoint(point)
	}
}

func SendResponses() {
	// Write batch
	err := DBClient.Write(respBP)
	if err != nil {
		glog.Errorf("Fail to write to db: %s", err.Error())
	}
	// Create new batch
	respBP, err = client_v2.NewBatchPoints(client_v2.BatchPointsConfig{
		Database:  "spy",
		Precision: "ms",
	})
	if err != nil {
		glog.Fatalf("Fail to create points batch: %s", err)
		glog.Flush()
		panic(err)
	}
}

func AddPingResult(serviceName,namespace string, chaos *Chaos,podName,delay,loss string){
	// Create map
	tags := make(map[string]string)
	fields := make(map[string]interface{})

	tags["serviceName"] = serviceName
	tags["namespace"]=namespace
	tags["podName"]=podName
	fields["delay"] = delay
	fields["loss"] = loss

	if chaos == nil {
		fields["chaos-ingress"] = "none"
		fields["chaos-egress"] = "none"
		fields["chaos-replica"] = "none"
	} else {
		fields["chaos-ingress"] = chaos.Ingress
		fields["chaos-egress"] = chaos.Egress
		if chaos.Replica == 0 {
			fields["chaos-replica"] = "none"
		} else {
			fields["chaos-replica"] = strconv.Itoa(chaos.Replica)
		}
	}

	// Create point
	point, err := client_v2.NewPoint(
		"ping",
		tags,
		fields,
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
	err = DBClient.Write(pingBP)
	if err != nil {
		glog.Errorf("Fail to write to db: %s", err.Error())
	}
	// Create new batch
	pingBP, err = client_v2.NewBatchPoints(client_v2.BatchPointsConfig{
		Database:  "spy",
		Precision: "ms",
	})
	if err != nil {
		glog.Fatalf("Fail to create points batch: %s", err)
		glog.Flush()
		panic(err)
	}
}
