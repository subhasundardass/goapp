package core

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const bannerInner = 60 // inner width between │ and │

// PrintBanner prints a startup banner to stdout.
// Call this after Bootstrap and before Run in main.go:
//
//	core.PrintBanner(app.Server, cfg, len(app.Registry.Modules()))
func PrintBanner(server *fiber.App, cfg *Config, moduleCount int) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	border := strings.Repeat("─", bannerInner)
	divider := "├" + border + "┤"

	lines := []string{
		"┌" + border + "┐",
		centerLine(cfg.AppName),
		centerLine("Version : " + cfg.AppVersion),
		divider,
		rowLine("Environment ..........", cfg.AppEnv),
		rowLine("URL ..................", cfg.AppURL),
		rowLine("Database .............", cfg.DBDriver+"  "+dbPing(cfg)),
		rowLine("Auto Migration .......", boolLabel(cfg.AutoMigration)),
		divider,
		rowLine("Handlers .............", fmt.Sprintf("%d", server.HandlersCount())),
		rowLine("Modules ..............", fmt.Sprintf("%d", moduleCount)),
		rowLine("Prefork ..............", "Disabled"),
		rowLine("PID ..................", fmt.Sprintf("%d", os.Getpid())),
		divider,
		rowLine("Memory ...............", fmt.Sprintf("%.1f MB", float64(mem.Alloc)/1024/1024)),
		rowLine("Go runtime ...........", runtime.Version()),
		rowLine("Debug ................", boolLabel(cfg.Debug)),
		"└" + border + "┘",
	}

	fmt.Println()
	for _, l := range lines {
		fmt.Println(l)
	}
	fmt.Println()
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// centerLine returns a banner line with content centered inside │ │.
func centerLine(s string) string {
	s = truncate(s, bannerInner)
	total := bannerInner - len([]rune(s))
	left := total / 2
	right := total - left
	return "│" + strings.Repeat(" ", left) + s + strings.Repeat(" ", right) + "│"
}

// rowLine returns a banner line with label and value left-aligned inside │ │.
func rowLine(label, value string) string {
	content := label + " " + value
	return "│ " + padRight(content, bannerInner-1) + "│"
}

// padRight pads or truncates s to exactly width runes.
func padRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) > width {
		runes = runes[:width]
	}
	for len(runes) < width {
		runes = append(runes, ' ')
	}
	return string(runes)
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) > max {
		return string(r[:max])
	}
	return s
}

func boolLabel(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// dbPing does a lightweight check on the DSN config.
// Replace with a real db.Ping() if you pass *ent.Client into PrintBanner later.
func dbPing(cfg *Config) string {
	if cfg.DBDSN == "" {
		return "not configured"
	}
	return "✓ connected"
}
