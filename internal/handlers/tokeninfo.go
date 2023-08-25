package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// TokenInfo godoc
// @Summary      TokenInfo endpoint
// @Description  return info of the given token
// @ID           get-token-info
// @Tags         TokenInfo
// @Produce      json
// @Success      200  {string}  string
// @Router       /tokeninfo [get]
func tokeninfo(c *fiber.Ctx) error {
	var scopes = c.Locals("scopes")
	var client_id = c.Locals("client_id")
	return c.JSON([]interface{}{scopes, client_id})
}
