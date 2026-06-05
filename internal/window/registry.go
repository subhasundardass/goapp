package window

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

type Factory func(payload any) templ.Component

var registry = map[string]Factory{}

func Register(key string, factory Factory) {
	registry[key] = factory
}

func Resolve(key string, payload any) templ.Component {
	if f, ok := registry[key]; ok {
		return f(payload)
	}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := w.Write([]byte("window not found"))
		return err
	})
}
