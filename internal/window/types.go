package window

type Size string

const (
	Sm Size = "sm"
	Md Size = "md"
	Lg Size = "lg"
)

type Window struct {
	Key     string `json:"key"`
	Title   string `json:"title"`
	Size    Size   `json:"size"`
	Payload any    `json:"payload,omitempty"`
}

type Command struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

type Response struct {
	Commands []Command `json:"commands"`
}
