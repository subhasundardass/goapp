.PHONY: run build templ modules dev css css-watch wasm clean module

CSS_IN  = web/static/css/input.css
CSS_OUT = web/static/css/app.css

# ── Scaffold a new module ─────────────────────────────────────────────────────
# Usage: make module name=invoice
#        make module name=sales_order
module:
	@if [ -z "$(name)" ]; then \
		echo "usage: make module name=<module_name>"; \
		exit 1; \
	fi
	go run ./cmd/gen/scaffold --name=$(name)
	templ generate
	$(MAKE) modules
	go mod tidy
	@echo "✅ module '$(name)' ready"

# ── Generate module imports ───────────────────────────────────────────────────
modules:
	go run ./cmd/gen/modules
	@echo "✅ modules_gen.go updated"

# ── Generate templ components ─────────────────────────────────────────────────
templ:
	templ generate
	@echo "✅ templ generated"

# ── Tailwind CSS ──────────────────────────────────────────────────────────────
css:
	npm run tailwind -- -i $(CSS_IN) -o $(CSS_OUT) --minify
	@echo "✅ css built"

css-watch:
	npm run tailwind -- -i $(CSS_IN) -o $(CSS_OUT) --watch

# ── Run (one shot, no watching) ───────────────────────────────────────────────
run: modules templ
	go run ./cmd/web

# ── Build ─────────────────────────────────────────────────────────────────────
build: modules templ css
	@mkdir -p bin
	go build -o bin/goapp ./cmd/web
	@echo "✅ built bin/goapp"

# ── Watch (hot reload) ────────────────────────────────────────────────────────
dev: modules templ css
	air

# ── WASM build ────────────────────────────────────────────────────────────────
wasm:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/web
	@echo "✅ wasm built"

# ── Clean ─────────────────────────────────────────────────────────────────────
clean:
	rm -rf bin/ tmp/ web/app.wasm $(CSS_OUT)
	@echo "✅ cleaned"

# =========================================================
# Kill Port
# =========================================================
kill:
	@echo "Killing process on port 8000..."
	@kill -9 $$(lsof -t -i:8000) 2>/dev/null || true
	@echo "Done"	