package server

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"            // Echo core library
	"github.com/labstack/echo/v4/middleware" // Optional Echo middleware
	_ "github.com/lukegriffith/SSHTrust/docs"
	"github.com/lukegriffith/SSHTrust/pkg/auth"
	"github.com/lukegriffith/SSHTrust/pkg/certStore"
	"github.com/lukegriffith/SSHTrust/pkg/handlers" // Import your cert package
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	Port = ":8080"
)

// generateRandomJWTSecret generates a cryptographically secure random JWT secret
func generateRandomJWTSecret() []byte {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate random JWT secret:", err)
	}
	log.Println("Generated random JWT secret for this session")
	return bytes
}

// loadJWTSecret loads JWT secret from environment variable or generates a random one
func loadJWTSecret() []byte {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		// If JWT_SECRET is base64 encoded, decode it
		if decoded, err := base64.StdEncoding.DecodeString(secret); err == nil && len(decoded) >= 32 {
			log.Println("Using JWT secret from JWT_SECRET environment variable (base64 decoded)")
			return decoded
		}
		// Otherwise use as-is if it's long enough
		if len(secret) >= 32 {
			log.Println("Using JWT secret from JWT_SECRET environment variable")
			return []byte(secret)
		}
		log.Println("JWT_SECRET environment variable is too short (minimum 32 characters), generating random secret")
	}
	return generateRandomJWTSecret()
}

// SetupServer configures the Echo instance and returns it for testing or running
func SetupServer(noAuth bool) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(2)

	// Optional middleware for logging and recovery
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// TODO: do this properly and salt
	auth.Users = &certStore.InMemoryUserList{}
	// Load JWT secret from environment or generate random
	auth.JWTSecret = loadJWTSecret()

	// Serve the Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	App := handlers.App{
		Store: certStore.NewInMemoryCaStore(),
	}

	var ca *echo.Group
	e.POST("/login", auth.Login)
	e.POST("/register", auth.Register)
	if noAuth {
		ca = e.Group("/CA")
	} else {
		ca = e.Group("/CA", echojwt.WithConfig(echojwt.Config{
			SigningKey: auth.JWTSecret,
		}))
	}
	// Define routes and their corresponding handlers
	ca.GET("", App.ListCA)         // List CAs
	ca.POST("", App.CreateCA)      // Create a new CA
	ca.GET("/:id", App.GetCA)      // Get a specific CA by ID
	ca.POST("/:id/Sign", App.Sign) // Sign a public key with a specific CA
	return e
}
