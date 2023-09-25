package middlewares

import (
	"fmt"

	spotifyUC "github.com/chunnior/spotify/internal/usecase/spotify"
	"github.com/gofiber/fiber/v2"
)

func ValidateState(spotifyUseCase *spotifyUC.Implementation) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		state := c.Query("state")
		stateKey := fmt.Sprintf(`state-%s`, state)
		ctx := c.Context()
		_, value := spotifyUseCase.CacheStates.Get(ctx, stateKey)
		if state != value {
			return c.Status(fiber.StatusUnauthorized).SendString("Error: Ivalid state")
		}
		go spotifyUseCase.CacheStates.Delete(ctx, stateKey)
		return c.Next()
	}
}
