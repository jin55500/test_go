package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "test/go/docs"

	AccountingController "test/go/controller/accounting"
	AuthController "test/go/controller/auth"
	UserController "test/go/controller/user"
	"test/go/database"
	"test/go/middleware"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	godotenv.Load(".env")

	database.ConnectDb()
	router := gin.Default()
	router.Use(cors.Default())

	routerUser := router.Group("/user")
	{
		routerUser.POST("/register", AuthController.RegisterUser)
		routerUser.POST("/login", AuthController.Login)
		routerUser.GET("/me", middleware.ValidateToken(), UserController.GetUser)
		routerUser.PATCH("/me", middleware.ValidateToken(), UserController.UpdateUser)
	}

	routerAccounting := router.Group("/accounting")
	{
		routerAccounting.POST("/transfer", middleware.ValidateToken(), AccountingController.Transfer)
		routerAccounting.GET("/transfer-list", middleware.ValidateToken(), AccountingController.GetTransferList)
	}

	docs.SwaggerInfo.Title = "Swagger Example API"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type ping struct {
	Message string `json:"message" example:"pong" `
}
