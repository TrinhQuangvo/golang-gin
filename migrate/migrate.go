package main

import (
	"go-crud/initilizers"
	"go-crud/models"
)

func init() {
	initilizers.LoadEnvVariables()
	initilizers.ConnectToDatabase()
}
func main() {
	initilizers.DB.AutoMigrate(
		&models.Auth{},
		&models.Profile{},
		&models.AuthRoles{},
		&models.Role{},
		&models.Post{},
	)
}
