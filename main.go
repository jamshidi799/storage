package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"storage/docs"
	"storage/record"
	"storage/user"
	"time"
)

func main() {
	if err := Run(); err != nil {
		//log.Fatal(err)
	}

	time.Sleep(time.Hour)
}

func Run() error {
	loadEnv()
	postgresDB, err := initPostgresDB()
	if err != nil {
		return err
	}

	r := gin.Default()
	environment := os.Getenv("ENVIRONMENT")
	if environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	basePath := "/api"
	docs.SwaggerInfo.BasePath = basePath
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group(basePath)

	uGroup := api.Group("user")
	uRepo := user.NewPostgresUserRepository(postgresDB)

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtTokenGenerator := user.NewJwtTokenGenerator(jwtSecret)

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
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DATABASE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, dbUser, password, database)

	mode := os.Getenv("ENVIRONMENT")

	var config gorm.Config
	if mode == "dev" {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	return gorm.Open(postgres.Open(dsn), &config)
}

func loadEnv() {
	env := os.Getenv("POSTGRES_HOST")
	println("env: ", env)
	if env == "" {
		_ = godotenv.Load()
	}
}
