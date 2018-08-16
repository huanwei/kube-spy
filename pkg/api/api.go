package api

import (
	"github.com/gorilla/mux"
	client_v2 "github.com/influxdata/influxdb/client/v2"
)

type ApiServer struct {
	Router   *mux.Router
	InfluxDB *client_v2.Client
}

// See https://github.com/influxdata/influxdb/tree/master/client

func (a *ApiServer) Initialize(user, password, dbname string) {
	client_v2.NewHTTPClient(client_v2.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: user,
		Password: password,
	})
}

func (a *ApiServer) Run(addr string) {}
