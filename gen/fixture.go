package gen

import "github.com/kettek/morogue/id"

type Fixture struct {
	ID   id.UUID
	Keys map[string]id.UUID
	Rows []string
}
