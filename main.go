package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"storage/user"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	postgresDB, err := initPostgresDB()
	if err != nil {
		return err
	}

	r := gin.Default()

	basePath := "api"
	api := r.Group(basePath)

	uGroup := api.Group("user")
	uRepo := user.NewPostgresUserRepository(postgresDB)
	uService := user.NewUserService(uRepo)
	user.NewUserController(uGroup, uService)

	return r.Run()
}

func initPostgresDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=storage port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
