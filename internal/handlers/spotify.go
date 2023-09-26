package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/chunnior/spotify/internal/models"
	spotifyUC "github.com/chunnior/spotify/internal/usecase/spotify"
	"github.com/chunnior/spotify/internal/util/log"
	"github.com/chunnior/spotify/pkg/tracing"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	useCase spotifyUC.UseCase
}

func NewHandler(useCase spotifyUC.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Save(c *fiber.Ctx) error {
	ctx := startTracing(c, "saveTaste")
	defer endTracing(ctx)

	log.InfofWithContext(ctx, "[TASTE_POST] Creating Taste")

	var req models.DataRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return respondWithError(c, http.StatusBadRequest, "invalid_body")
	}

	taste, err := h.useCase.Save(ctx, req)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "internal_error")
	}

	log.InfofWithContext(ctx, "[TASTE_POST] Taste created: %s", taste.Data)
	return c.Status(http.StatusOK).JSON(fiber.Map{"data": taste})
}

func (h *Handler) Get(c *fiber.Ctx) error {
	ctx := startTracing(c, "getData")
	defer endTracing(ctx)

	dataType, userId := c.Params("dataType"), c.Params("userId")
	if dataType == "" || userId == "" {
		return respondWithError(c, http.StatusBadRequest, "invalid_params")
	}

	log.InfofWithContext(ctx, "[GETDATA] Getting data %s for user %s", dataType, userId)

	res, err := h.useCase.Get(ctx, dataType, userId)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}
	if res == nil {
		return respondWithError(c, http.StatusNotFound, "data not found")
	}

	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	ctx := startTracing(c, "login")
	defer endTracing(ctx)

	authURL, state, err := h.useCase.Login(ctx)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}
	if authURL == nil {
		return respondWithError(c, http.StatusInternalServerError, "Cannot get Spotify auth url")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "oauth_state",
		Value: state,
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{"url": authURL})
}

func (h *Handler) Callback(c *fiber.Ctx) error {
	ctx := startTracing(c, "callback")
	defer endTracing(ctx)

	code := c.Query("code")
	if code == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	state := c.Query("state")
	if state == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := h.useCase.HandleCallback(ctx, code, state)

	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusOK).JSON(user)
}

func startTracing(c *fiber.Ctx, segmentName string) context.Context {
	ctx := tracing.CreateContextWithTransaction(c)
	txn := tracing.GetTransactionFromContext(ctx)
	txn.StartSegment(segmentName)
	return ctx
}

func endTracing(ctx context.Context) {
	txn := tracing.GetTransactionFromContext(ctx)
	txn.End()
}

func respondWithError(c *fiber.Ctx, status int, msg string) error {
	log.ErrorfWithContext(c.Context(), msg)
	return c.Status(status).JSON(formatResponse(status, msg))
}

func formatResponse(status int, msg string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": msg,
	}
}
