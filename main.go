package main

import (
	"fmt"
	"simple-api-example/database"
	"simple-api-example/router"
	"simple-api-example/utils"
)

func main() {
	fmt.Println("hello world")
	utils.LoadEnv()

	database.SetDB()
	// models.CreateAllTablesIfNotExists(database.DB)
	router := router.SetupRouter()
	router.Run(":8080")
}
