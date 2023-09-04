package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chunnior/spotify/internal/models"
	"github.com/chunnior/spotify/internal/usecase/spotify"
	"github.com/chunnior/spotify/internal/util/log"
	"github.com/chunnior/spotify/internal/util/tracing"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	useCase spotify.UseCase
}

func NewHandler(useCase spotify.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// Save is a function that saves a taste
func (h *Handler) Save(c *fiber.Ctx) error {
	ctx := tracing.CreateContextWithTransaction(c)
	txn := tracing.GetTransactionFromContext(ctx)
	segment := txn.StartSegment("saveTaste")
	defer segment.End()

	log.InfofWithContext(ctx, "[TASTE_POST] Creating Taste")

	var req models.TasteRequest

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		log.ErrorfWithContext(ctx, "[TASTE_POST] Error unmarshalling request: %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid_body"})
	}

	//TODO: AGREGAR VALIDACION en MODELS
	// err = req.Validate()
	// if err != nil {
	// 	log.ErrorfWithContext(ctx, "[TASTE_POST] Error validating request: %s", err.Error())
	// 	return c.Status(http.StatusBadRequest).JSON((fiber.Map{"message": err.Error()}))
	// }

	taste, err := h.useCase.Save(ctx, req)

	if err != nil {
		log.ErrorfWithContext(ctx, "[TASTE_POST] Error saving taste: %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "internal_error"})
	}

	log.InfofWithContext(ctx, "[TASTE_POST] Taste created: %s", taste.ID)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": taste,
	})
}

func (h *Handler) Get(c *fiber.Ctx) error {
	ctx := tracing.CreateContextWithTransaction(c)
	txn := tracing.GetTransactionFromContext(ctx)
	segment := txn.StartSegment("getData")
	defer segment.End()

	dataType := c.Params("dataType")
	userId := c.Params("userId")

	if dataType == "" || userId == "" {
		log.ErrorfWithContext(ctx, "[GETDATA] Error getting data: %s", "invalid_params")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid_params"})
	}

	log.InfofWithContext(ctx, "[GETDATA] Getting data %s for user %s", dataType, userId)

	res, err := h.useCase.Get(ctx, dataType, userId)
	if err != nil {
		log.Errorf("Error getting data: %s", err.Error())
		return c.Status(500).JSON(formatResponse(500, err.Error()))
	}
	if res == nil {
		return c.Status(404).JSON(formatResponse(404, "data not found"))
	}

	return c.Status(200).JSON(res)

}

func formatResponse(status int, msg string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": msg,
	}
}
