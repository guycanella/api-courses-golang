package httpx

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

var debug bool

func SetDebug(b bool) { debug = b }

func InternalServerError(c *fiber.Ctx, err error) error {
	if err != nil {
		log.Printf("500 internal error: %v", err)
	}

	msg := "internal server error"
	if debug && err != nil {
		msg = err.Error()
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
}
