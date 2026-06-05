package handler

import (
	"strconv"

	"goapp/internal/core"
	"goapp/internal/form"
	"goapp/modules/accounting/model"
	"goapp/modules/accounting/pages"
	"goapp/modules/accounting/service"

	"github.com/gofiber/fiber/v2"
)

type AccountingHandler struct {
	service *service.Services
}

func NewAccountingHandler(svc *service.Services) *AccountingHandler {
	return &AccountingHandler{service: svc}
}

// Index renders the landing page.
func (h *AccountingHandler) Index(c *fiber.Ctx) error {
	return core.Render(c, pages.AccountingPage("Accounting"))
}

// *** Group
func (h *AccountingHandler) GroupList(c *fiber.Ctx) error {
	params := core.PaginationParams{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
		Search:  c.Query("q"),
	}
	result, err := h.service.Accounting.GroupListPaginated(c.Context(), params)
	if err != nil {
		return core.InternalError(c, err)
	}
	return core.Render(c, pages.GroupList("Group Master", result))
}

// *** Ledger
func (h *AccountingHandler) LedgerList(c *fiber.Ctx) error {
	params := core.PaginationParams{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
		Search:  c.Query("q"),
	}
	result, err := h.service.Accounting.LedgerListPaginated(c.Context(), params)
	if err != nil {
		return core.InternalError(c, err)
	}
	return core.Render(c, pages.LedgerList("Ledger Master", result))
}

// *** Party
func (h *AccountingHandler) PartyList(c *fiber.Ctx) error {
	params := core.PaginationParams{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
		Search:  c.Query("q"),
	}
	result, err := h.service.Accounting.PartyListPaginated(c.Context(), params)
	if err != nil {
		return core.InternalError(c, err)
	}
	return core.Render(c, pages.PartyList("Party Master", result))
}

func (h *AccountingHandler) PartyNew(c *fiber.Ctx) error {

	var partyForm = model.PartyFormModel()
	engine, err := form.NewEngine(partyForm)
	if err != nil {
		return core.InternalError(c, err)
	}

	// Country
	countries, err := h.service.Accounting.CountryList(c.Context())
	if err != nil {
		return err
	}
	countryList := make([]form.Option, 0, len(countries))
	for _, l := range countries {

		countryList = append(countryList, form.Option{
			Value:    strconv.Itoa(l.ID),
			Label:    l.Name,
			Disabled: false,
		})
	}
	fieldCountry := partyForm["countries"]
	fieldCountry.Options = countryList
	partyForm["ledgers"] = fieldCountry

	// Ledgers
	ledgers, err := h.service.Accounting.LedgerList(c.Context())
	if err != nil {
		return err
	}
	lopts := make([]form.Option, 0, len(ledgers))
	for _, l := range ledgers {

		lopts = append(lopts, form.Option{
			Value:    strconv.Itoa(l.ID),
			Label:    l.Name,
			Disabled: false,
		})
	}

	fieldL := partyForm["ledgers"]
	fieldL.Options = lopts
	partyForm["ledgers"] = fieldL

	// _ = engine.SetField("ledgers", lopts)

	ui := form.ProjectMap(partyForm, engine.GetState())

	return core.Render(c, pages.PartyNew("New Party", ui))
}

// func (h *AccountingHandler) LedgerEdit(c *fiber.Ctx) error {
// 	id, _ := strconv.Atoi(c.Params("id"))

// 	ledger, err := h.service.Accounting.FindLedgerByID(c.Context(), id)
// 	if err != nil {
// 		return core.NotFound(c, "ledger")
// 	}

// 	window.Open(c, "ledger.form", window.Props{
// 		Title:   "Edit Ledger",
// 		Size:    window.SizeLg,
// 		OnClose: "#ledger-list",
// 		Payload: formwindow.FormPayload{
// 			SubmitURL:   fmt.Sprintf("/accounting/ledger/%d/update", id),
// 			SubmitLabel: "Update",
// 			InitialValues: map[string]any{
// 				"name": ledger,
// 			},
// 		},
// 	})

// 	return window.Render(c, window.FromCtx(c))
// }

// func (h *AccountingHandler) LedgerEdit(c *fiber.Ctx) error {
// 	id, _ := strconv.Atoi(c.Params("id"))

// 	ledger, err := h.service.Accounting.FindLedgerByID(c.Context(), id)
// 	if err != nil {
// 		return core.NotFound(c, "ledger")
// 	}

// 	ui := window.New().
// 		Open(window.Window{
// 			Key:   "ledger.form",
// 			Title: "Edit Ledger",
// 			Size:  window.Lg,
// 			Payload: formwindow.FormPayload{
// 				SubmitURL: fmt.Sprintf("/accounting/ledger/%d/update", id),
// 				InitialValues: map[string]any{
// 					"name": ledger,
// 				},
// 			},
// 		}).
// 		Toast("Ledger loaded")

// 	return c.JSON(ui.Build())
// }

// func (h *AccountingHandler) LedgerEdit(c *fiber.Ctx) error {
// 	id, _ := strconv.Atoi(c.Params("id"))

// 	// ledger, err := h.service.Accounting.FindLedgerByID(c.Context(), id)
// 	// if err != nil {
// 	// 	return core.NotFound(c, "ledger")
// 	// }

// 	ui := window.New().
// 		Open(window.Window{
// 			Key:   "ledger.form",
// 			Title: "Edit Ledger",
// 			Size:  window.Lg,
// 			Payload: formwindow.FormPayload{
// 				SubmitURL:   fmt.Sprintf("/accounting/ledger/%d/update", id),
// 				SubmitLabel: "Update",
// 				InitialValues: map[string]any{
// 					"name": "Testing",
// 				},
// 			},
// 		}).
// 		Toast("Ledger loaded")

// 	return c.JSON(ui.Build())
// }
