package window

type Engine struct {
	cmds []Command
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) Open(w Window) *Engine {
	e.cmds = append(e.cmds, Command{
		Type: "open",
		Data: w,
	})
	return e
}

func (e *Engine) Toast(msg string) *Engine {
	e.cmds = append(e.cmds, Command{
		Type: "toast",
		Data: map[string]any{
			"message": msg,
		},
	})
	return e
}

func (e *Engine) Close() *Engine {
	e.cmds = append(e.cmds, Command{
		Type: "close",
	})
	return e
}

func (e *Engine) Build() Response {
	return Response{Commands: e.cmds}
}
