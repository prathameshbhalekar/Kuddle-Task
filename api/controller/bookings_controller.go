package controller

import (
	"time"

	"github.com/KuddleTask/api/db"
	"github.com/KuddleTask/api/schemas"
	"github.com/gofiber/fiber/v2"
)

func BookSlot(ctx *fiber.Ctx) error {

	body := new(schemas.Registration)

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	var class schemas.Class

	//find class
	if err := db.GetDB().
		Table("classes").
		Where("CLASS_UUID = ?", body.ClassUuid).
		First(&class).
		Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	body.BookedAt = int(time.Now().Unix())
	body.IsCancelled = false

	// check for full class
	if class.Members >= class.Capacity {
		body.IsWaiting = true

		if err := db.GetDB().
			Table("registrations").
			Create(body).
			Error; err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"success": false,
				"msg":     err.Error(),
			})
		}

		return ctx.Status(400).JSON(fiber.Map{
			"success": true,
			"msg":     "Added to waiting list",
			"result":  body,
		})
	}

	body.IsWaiting = false

	if err := db.GetDB().
		Table("registrations").
		Create(body).
		Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	if err := db.GetDB().
		Exec(`UPDATE CLASSES SET MEMBERS = ? WHERE CLASS_UUID = ?`, class.Members+1, class.ClassUuid).
		Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	return ctx.Status(400).JSON(fiber.Map{
		"success": true,
		"msg":     "Added to waiting list",
		"result":  body,
	})

}

func CancelSlot(ctx *fiber.Ctx) error {

	body := new(schemas.Registration)

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	if err := db.GetDB().
		Exec(`UPDATE registrations SET IS_CANCELLED = TRUE WHERE CLASS_UUID = ? AND USER_UUID = ?`, body.ClassUuid, body.UserUuid).
		Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	var new_reg schemas.Registration
	if err := db.GetDB().
		Table("registrations").
		Where("CLASS_UUID = ? AND IS_WAITING", body.ClassUuid).
		Order("BOOKED_AT ASC").
		First(&new_reg).
		Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	if err := db.GetDB().
		Exec(`UPDATE registrations SET IS_WAITING = FALSE WHERE CLASS_UUID = ? AND USER_UUID = ?`, new_reg.ClassUuid, new_reg.UserUuid).
		Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	return ctx.Status(400).JSON(fiber.Map{
		"success": true,
		"msg":     "Cancelled Successfully",
		"result":  body,
	})

}
