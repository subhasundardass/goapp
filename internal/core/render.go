package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

// Render writes a templ component to the Fiber response as text/html.
func Render(c *fiber.Ctx, component templ.Component, options ...func(*RenderOptions)) error {
	opts := defaultRenderOptions()
	for _, o := range options {
		o(opts)
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Status(opts.StatusCode)

	return component.Render(c.Context(), c.Response().BodyWriter())
}

// RenderSSE writes a templ component as a Datastar SSE merge-fragments event.
// Use this inside a streaming handler to push HTML patches to the browser.
func MergeSignalsSSE(c *fiber.Ctx, signals map[string]any) error {
	data, err := json.Marshal(signals)
	if err != nil {
		return fmt.Errorf("sse signals marshal: %w", err)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Response().ImmediateHeaderFlush = true

	payload := fmt.Sprintf("event: datastar-merge-signals\ndata: signals %s\n\n", string(data))
	fmt.Println("SSE PAYLOAD:", payload) // ← check terminal

	fmt.Fprint(c.Response().BodyWriter(), payload)
	return nil
}

// StreamSSE sets up the response headers for a long-lived SSE connection
// and calls the provided handler function. The handler should call RenderSSE
// repeatedly to push updates.
//
// Example:
//
//	func (h *Handler) Stream(c *fiber.Ctx) error {
//	    return core.StreamSSE(c, func() error {
//	        items, _ := h.service.List()
//	        return core.RenderSSE(c, components.Table(items))
//	    })
//	}
func StreamSSE(c *fiber.Ctx, fn func() error) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)

	return fn()
}

// ── Options ───────────────────────────────────────────────────────────────────

type RenderOptions struct {
	StatusCode int
}

func defaultRenderOptions() *RenderOptions {
	return &RenderOptions{StatusCode: http.StatusOK}
}

// WithStatus sets a custom HTTP status code on the rendered response.
func WithStatus(code int) func(*RenderOptions) {
	return func(o *RenderOptions) {
		o.StatusCode = code
	}
}
