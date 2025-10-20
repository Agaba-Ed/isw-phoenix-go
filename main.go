package main

import (
	"log"
	"os"

	"isw-phoenix-go/internal/isw"
	"isw-phoenix-go/internal/logger"
	"isw-phoenix-go/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	logger.Log.Info("Starting ISW Phoenix Server")

	// Load environment variables
	if err := godotenv.Load("config/.env"); err != nil {
		logger.Log.WithError(err).Warn("Environment file not found")
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create configuration
	config := &isw.Configuration{
		APIURL:           getEnv("API_URL", ""),
		BillersAPIURL:    getEnv("BILLERS_API_URL", ""),
		ClientID:         getEnv("CLIENT_ID", ""),
		OwnerPhoneNumber: getEnv("OWNER_PHONE_NUMBER", ""),
		Passphrase:       getEnv("PASSPHRASE", ""),
		ClientSecretKey:  getEnv("CLIENT_SECRET_KEY", ""),
		TerminalID:       getEnv("TERMINAL_ID", ""),
		AppVersion:       getEnv("APP_VERSION", ""),
		SerialID:         getEnv("SERIAL_ID", getEnv("serialId", "")),
	}

	// Create ISW client
	iswClient := isw.NewClient(config)

	// Create handlers
	handlers := server.NewHandlers(iswClient)

	// Setup Gin router
	router := gin.Default()
	router.Use(server.RequestLoggerMiddleware())
	router.Use(gin.Recovery())

	// Setup routes
	setupRoutes(router, handlers)

	// Start server
	port := getEnv("PORT", "3000")
	logger.Log.WithField("port", port).Info("Server starting")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Log.WithError(err).Fatal("Failed to start server")
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine, handlers *server.Handlers) {
	// API routes
	router.POST("/validateCustomer", handlers.ValidateCustomer)
	router.GET("/getbillerbycategory/:id", handlers.GetBillerCategories)
	router.GET("/getcategories", handlers.GetCategories)
	router.GET("/billeritems/:id", handlers.GetPaymentItems)
	router.POST("/payment", handlers.Payment)
	router.GET("/transStatus/:id", handlers.TransStatus)
	router.GET("/accountBalance", handlers.AccountBalance)
	router.GET("/keyPair", handlers.GenerateRSAKeyPair)
	router.POST("/clientRegistration", handlers.ClientRegistration)
	router.POST("/doKeyExchange", handlers.DoKeyExchange)
	router.GET("/generateECDHKeyPair", handlers.GenerateECDHKeyPair)
	router.GET("/util", handlers.Util)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
