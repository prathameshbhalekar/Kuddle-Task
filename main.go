package main

import (
	"fmt"
	"strings"

	"github.com/KuddleTask/api/db"
	"github.com/KuddleTask/api/router"
	"github.com/KuddleTask/api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func main() {

	//Set global configuration
	utils.ImportEnv()

	// Create Fiber
	app := fiber.New(fiber.Config{})

	app.Get("/", healthCheck)
	app.Get("/health", healthCheck)

	app.Use(logger.New(logger.Config{Next: func(c *fiber.Ctx) bool {
		return strings.HasPrefix(c.Path(), "api")
	}}))

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "",
		AllowHeaders: "*",
	}))

	// Mount Routes
	router.MountPublicRoutes(app)

	// Get Port
	port := utils.GetPort()

	db.Connect()

	// Start Fiber
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

}
