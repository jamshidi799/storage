package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"storage/record"
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
	// todo: set gin mode

	basePath := "api"
	api := r.Group(basePath)

	uGroup := api.Group("user")
	uRepo := user.NewPostgresUserRepository(postgresDB)
	jwtTokenGenerator := user.NewJwtTokenGenerator("secret")
	uService := user.NewUserService(uRepo, jwtTokenGenerator)
	uHandler := user.NewUserController(uGroup, uService)

	rGroup := api.Group("record")
	rGroup.Use(uHandler.JwtAuthMiddleware())
	rRepo := record.NewPostgresRecordRepository(postgresDB)
	rService := record.NewRecordService(rRepo)
	record.NewRecordController(rGroup, rService)

	return r.Run()
}

func initPostgresDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=storage port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		// todo: enable verbose mode in development environment
		Logger: logger.Default.LogMode(logger.Info),
	})
}
