package main

import (
	"os"
	"simple-api-example/database"
	"simple-api-example/router"
	"simple-api-example/utils"
)

func main() {
	if os.Getenv("USE_K8S") == "" {
		utils.LoadEnv()
	}

	database.SetDB()
	// models.CreateAllTablesIfNotExists(database.DB)
	router := router.SetupRouter()
	router.Run(":8080")
}
