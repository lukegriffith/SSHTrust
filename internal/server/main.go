package server

import (
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

// SetupServer configures the Echo instance and returns it for testing or running
func SetupServer(noAuth bool) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(2)

	// Optional middleware for logging and recovery
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// TODO: do this properly, externally, hashed + salted, etc
	auth.Users = &certStore.InMemoryUserList{}

	auth.JWTSecret = []byte("secret")

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
