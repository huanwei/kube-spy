package main

//api program entrypoint

import (
	"github.com/huanwei/kube-spy/pkg/api"
)

func main() {

	//todo: api server
	a := api.ApiServer{}
	a.Initialize("DB_USERNAME", "DB_PASSWORD", "spy")
	a.Run(":8080")


}
