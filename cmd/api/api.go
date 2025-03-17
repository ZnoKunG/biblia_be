package main

import (
	"biblia-be/internal/db"
	"biblia-be/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type application struct {
	config config
}

type config struct {
	host string
	addr string
	db   dbConfig
}

type dbConfig struct {
	user         string
	host         string
	password     string
	db_name      string
	db_addr      string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  int
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Testing purpose
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	userHandler := handler.UserHandler{}
	userHandler.Initialize(db)

	router.GET("user", userHandler.GetUsers)
	router.GET("user/:id", userHandler.GetUser)
	router.POST("user", userHandler.CreateUser)
	router.PUT("user/:id", userHandler.UpdateUser)
	router.DELETE("user/:id", userHandler.DeleteUser)

	recordHandler := handler.RecordHandler{}
	recordHandler.Initialize(db)

	router.GET("record", recordHandler.GetRecords)
	router.GET("record/:id", recordHandler.GetRecord)
	router.POST("record", recordHandler.CreateRecord)
	router.PUT("record/:id", recordHandler.UpdateRecord)
	router.DELETE("record/:id", recordHandler.DeleteRecord)

	return router
}

// @title Biblia Backend API
// @version 1.0
// @description This is a biblia backend API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host biblia.swagger.io:3000
// @BasePath /v2
func (app *application) run() {
	db, err := db.NewDB(
		app.config.db.host,
		app.config.db.user,
		app.config.db.password,
		app.config.db.db_name,
		app.config.db.db_addr,
		app.config.db.maxOpenConns,
		app.config.db.maxIdleConns,
		app.config.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	router := setupRouter(db)

	url := ginSwagger.URL("http://.swagger.io:3000/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run(app.config.addr)
}
