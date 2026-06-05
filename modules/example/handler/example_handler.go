package handler

import (
	"fmt"
	"strconv"

	"goapp/internal/core"
	"goapp/internal/form"
	"goapp/modules/example/model"
	"goapp/modules/example/pages"
	"goapp/modules/example/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

var exampleForm = model.ExampleFormModel()

// Index renders the landing page for this module.
func (h *Handler) Index(c *fiber.Ctx) error {

	// example := model.ExampleFormModel()
	engine, err := form.NewEngine(exampleForm)

	if err != nil {
		return core.InternalError(c, err)
	}

	// pass 2.0 (float64) not 2 (int) — compute functions use float64
	_ = engine.SetField("itemName", "Laptop")
	_ = engine.SetField("qty", 2.0)
	_ = engine.SetField("rate", 1000.0)
	_ = engine.SetField("taxPercent", 10.0)

	// Use ProjectMap when you need key-based access like ui["payment_mode"]
	ui := form.ProjectMap(exampleForm, engine.GetState())

	// For the template, pass the ordered slice so fields render in Order sequence
	// uiSlice := form.Project(example, engine.GetState())

	fmt.Printf("STATE: %+v\n", engine.GetState())
	return core.Render(c, pages.ExamplePage("Example", ui))
}

func (h *Handler) Compute(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return core.BadRequest(c, err)
	}

	engine, err := form.NewEngine(exampleForm)
	if err != nil {
		return core.InternalError(c, err)
	}

	for k, v := range body {
		_ = engine.SetField(k, v)
	}

	state := engine.GetState()

	fmt.Printf("BODY: %+v\n", body)
	fmt.Printf("STATE item_name: %v\n", state["itemName"])

	// Only push computed (Readable) fields back as signals
	return core.MergeSignalsSSE(c, map[string]any{
		"itemName":    state["itemName"],
		"amount":      state["amount"],
		"tax_amount":  state["tax_amount"],
		"grand_total": state["grand_total"],
	})
}

// List returns all records as JSON.
func (h *Handler) List(c *fiber.Ctx) error {
	params := core.PaginationParams{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
		Search:  c.Query("q"),
	}
	result, err := h.service.ListPaginated(c.Context(), params)
	if err != nil {
		return core.InternalError(c, err)
	}
	return core.OK(c, result)
}

// Show returns a single record by ID.
func (h *Handler) Show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return core.BadRequest(c, err)
	}
	item, err := h.service.FindByID(c.Context(), id)
	if err != nil {
		return core.NotFound(c, "example")
	}
	return core.OK(c, item)
}

// Create creates a new record.
func (h *Handler) Create(c *fiber.Ctx) error {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&body); err != nil {
		return core.BadRequest(c, err)
	}
	item, err := h.service.Create(c.Context(), body.Name, body.Description)
	if err != nil {
		return core.UnprocessableEntity(c, err)
	}
	return core.Created(c, item)
}

// Update modifies an existing record.
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return core.BadRequest(c, err)
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&body); err != nil {
		return core.BadRequest(c, err)
	}
	item, err := h.service.Update(c.Context(), id, body.Name, body.Description)
	if err != nil {
		return core.UnprocessableEntity(c, err)
	}
	return core.OK(c, item)
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
