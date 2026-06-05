package window

type FormPayload struct {
	InitialValues    map[string]any    `json:"initialValues"`
	SubmitURL        string            `json:"submitUrl"`
	SubmitLabel      string            `json:"submitLabel"`
	CancelURL        string            `json:"cancelUrl"`
	ValidationErrors map[string]string `json:"validationErrors,omitempty"`
}
