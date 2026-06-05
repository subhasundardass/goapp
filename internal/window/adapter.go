package window

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

type Adapter struct {
	schema any
}

// optional if you use schema system later
func NewAdapter(schema any) *Adapter {
	return &Adapter{schema: schema}
}

func (a *Adapter) Render(p FormPayload) templ.Component {

	// ⚠️ Replace this with your real templ form
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {

		_, err := fmt.Fprintf(w, `
<div class="p-4">
  <form action="%s" method="post">

    <label>Name</label>
    <input name="name" value="%v" />

    <button type="submit">%s</button>

  </form>
</div>
`,
			p.SubmitURL,
			p.InitialValues["name"],
			p.SubmitLabel,
		)

		return err
	})
}
