package app

import (
	"github.com/gorilla/mux"
	client_v2 "github.com/influxdata/influxdb/client/v2"
)

type App struct {
	Router *mux.Router
	InfluxDB     *client_v2.Client
}

// See https://github.com/influxdata/influxdb/tree/master/client

func (a *App) Initialize(user, password, dbname string) {}

func (a *App) Run(addr string) {}
