package main

import (
	"api-ticket/config"
	"api-ticket/internal/cmd"
	"api-ticket/internal/controllers/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title hepytic
// @version 1.0
// @description This is an API documentation for hepytic.
// @host localhost:8000
// @BasePath /v1
func main() {
	// for load godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// database connection
	db := config.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	//Gin
	app := gin.Default()
	app.UseRawPath = true
	app.UnescapePathValues = true
	app.RemoveExtraSlash = true

	//CORS Config
	corsMiddleware := func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}

	//Apply Middleware
	app.Use(corsMiddleware)

	//Rate Limiter, for what ?
	//var limiter = rate.NewLimiter(rate.Limit(10), 1) // Contoh: 10 permintaan per detik
	//app.Use(func(c *gin.Context) {
	//	if limiter.Allow() == false {
	//		c.JSON(http2.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
	//		c.Abort()
	//		return
	//	}
	//	c.Next()
	//})

	//Inject Service
	bannerSetup := cmd.InitBannerService(db)
	talentSetup := cmd.InitTalentService(db)
	eventSetup := cmd.InitEventService(db)
	authSetup := cmd.InitAuthService(db)
	http.RouteService(
		&app.RouterGroup,
		bannerSetup,
		talentSetup,
		eventSetup,
		authSetup,
	)

	app.Static("/public", "./public")

	//// router
	//r := routes.SetupRouter(db)
	app.Run(":8000")
}
