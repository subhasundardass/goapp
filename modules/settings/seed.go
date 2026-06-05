package settings

import (
	"context"
	"goapp/ent"
	"goapp/ent/setting"
)

var defaults = []struct{ key, value, label, group string }{
	{"app_name", "GoApp", "App Name", "general"},
	{"theme", "light", "Theme", "appearance"},
	{"app_url", "", "App URL", "general"},
}

func Seed(ctx context.Context, db *ent.Client) error {
	for _, d := range defaults {
		exists, err := db.Setting.Query().
			Where(setting.KeyEQ(d.key)).
			Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			continue // never overwrite existing values
		}
		if err := db.Setting.Create().
			SetKey(d.key).
			SetValue(d.value).
			SetLabel(d.label).
			SetGroup(d.group).
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}
