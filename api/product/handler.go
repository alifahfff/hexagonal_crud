package product

import (
	"github.com/gofiber/fiber/v2"
)

// ProductHandler  an interface with operations to be implemented by a specific handler, ie http, gRCP
type ProductHandler interface {
	Get(ctx *fiber.Ctx)
	Create(ctx *fiber.Ctx)
	Update(ctx *fiber.Ctx)
	Delete(ctx *fiber.Ctx)
	GetAll(ctx *fiber.Ctx)
}
