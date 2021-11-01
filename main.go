package main

import (
	"fmt"
	"simple-api-example/router"
	"simple-api-example/utils"
)

func main() {
	fmt.Println("hello world")
	utils.LoadEnv()

	router := router.SetupRouter()
	router.Run(":8080")
}
