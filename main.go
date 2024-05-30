package main

import (
	"go-crud/initilizers"
	"go-crud/router"
)

func init() {
	initilizers.LoadEnvVariables()
	initilizers.ConnectToDatabase()
}

func main() {
	r := router.SetupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080
}
