package handler

import (
	"strconv"

	"goapp/internal/core"
	"goapp/modules/settings/pages"
	"goapp/modules/settings/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

// Index renders the landing page for this module.
func (h *Handler) Index(c *fiber.Ctx) error {
	return core.Render(c, pages.SettingsPage("Settings"))
}

// Delete removes a record by ID.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return core.BadRequest(c, err)
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		return core.InternalError(c, err)
	}
	return core.NoContent(c)
}
