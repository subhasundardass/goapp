package core

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ── JSON responses ────────────────────────────────────────────────────────────

func OK(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

// ── Error responses ───────────────────────────────────────────────────────────

func BadRequest(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

func Unauthorized(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"error":   "unauthorized",
	})
}

func Forbidden(c *fiber.Ctx) error {
	return c.Status(http.StatusForbidden).JSON(fiber.Map{
		"success": false,
		"error":   "forbidden",
	})
}

func NotFound(c *fiber.Ctx, resource string) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"error":   resource + " not found",
	})
}

func UnprocessableEntity(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

func InternalError(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"error":   "internal server error", // never expose raw err in production
	})
}

// ── Validation error response ─────────────────────────────────────────────────

// ValidationErrors returns a structured map of field → message pairs.
//
// Usage:
//
//	return core.ValidationErrors(c, map[string]string{
//	    "email": "must be a valid email address",
//	    "name":  "required",
//	})
func ValidationErrors(c *fiber.Ctx, fields map[string]string) error {
	return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
		"success": false,
		"errors":  fields,
	})
}

// ── Redirect ──────────────────────────────────────────────────────────────────

func Redirect(c *fiber.Ctx, url string) error {
	return c.Redirect(url, http.StatusSeeOther)
}
