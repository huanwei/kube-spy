package main

import "github.com/huanwei/kube-spy/app"

func main() {
	a := app.App{}
	a.Initialize("DB_USERNAME", "DB_PASSWORD", "rest_api_demo")
	a.Run(":8080")
}
