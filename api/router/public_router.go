package router

import (
	"github.com/KuddleTask/api/controller"
	"github.com/gofiber/fiber/v2"
)

func MountPublicRoutes(c *fiber.App) {
	api := c.Group("api")

	{
		api.Post("/booking/new", controller.BookSlot)
		api.Patch("/booking/cancel", controller.CancelSlot)
	}

}
