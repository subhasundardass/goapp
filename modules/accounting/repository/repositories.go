package repository

import "goapp/ent"

// Repositories groups all repositories for the Accounting module.
// Add new repository structs here as the module grows:
//
//	type Repositories struct {
//	    Accounting  *AccountingRepository
//	    Journal *JournalRepository
//	}
type Repositories struct {
	Accounting *AccountingRepository
}

// NewRepositories wires all repositories with the shared DB client.
func NewRepositories(db *ent.Client) *Repositories {
	return &Repositories{
		Accounting: NewAccountingRepository(db),
	}
}
