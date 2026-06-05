package example

import (
	"context"
	"goapp/ent"
)

// Seed inserts default records for the Example module.
// It is a no-op if records already exist.
// Called from PostRegister when APP_ENV=development.
func Seed(ctx context.Context, db *ent.Client) error {
	// TODO: replace with real ent query once schema is generated
	// count, err := db.Example.Query().Count(ctx)
	// if err != nil || count > 0 {
	//     return err
	// }
	// _, err = db.Example.Create().SetName("Sample Example").Save(ctx)
	// return err
	return nil
}
