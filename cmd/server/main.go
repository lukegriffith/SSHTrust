package main

import (
	"github.com/labstack/echo/v4"            // Echo core library
	"github.com/labstack/echo/v4/middleware" // Optional Echo middleware
	_ "github.com/lukegriffith/SSHTrust/docs"
	"github.com/lukegriffith/SSHTrust/pkg/certStore"
	"github.com/lukegriffith/SSHTrust/pkg/handlers" // Import your cert package
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	port = ":8080"
)

// SetupServer configures the Echo instance and returns it for testing or running
func SetupServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(2)

	// Optional middleware for logging and recovery
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve the Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	App := handlers.App{
		Store: certStore.NewInMemoryCaStore(),
	}

	// Define routes and their corresponding handlers
	e.GET("/CA", App.ListCA)         // List CAs
	e.POST("/CA", App.CreateCA)      // Create a new CA
	e.GET("/CA/:id", App.GetCA)      // Get a specific CA by ID
	e.POST("/CA/:id/Sign", App.Sign) // Sign a public key with a specific CA

	return e
}

func main() {
	e := SetupServer()

	e.Logger.Printf("SSHTrust Started on %s", port)
	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err)
	}
}
