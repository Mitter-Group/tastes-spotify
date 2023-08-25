package handlers

import "github.com/gofiber/fiber/v2"

// Liveness godoc
// @Summary      Liveness endpoint
// @Description  return 200 if the service its alive
// @ID           get-string-by-int
// @Tags         Healthcheck
// @Produce      json
// @Success      200  {string}  string
// @Router       /liveness [get]
func ping(c *fiber.Ctx) error {
	return c.SendString("Pong!")
}
