package main

import (
	"photo-app/app"
	"photo-app/database"
	_ "photo-app/docs"
	"photo-app/helpers"
)

//	@title						Photo App
//
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				JWT Bearer Token. Format: "Bearer <your-token-here>"
func main() {
	logger := helpers.NewLogger()
	config, err := helpers.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	db, err := database.New(config.DB)
	if err != nil {
		panic(err)
	}

	app := app.New(config.App, db, logger)

	if err := app.Start(); err != nil {
		panic(err)
	}
}
