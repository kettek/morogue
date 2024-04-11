package server

import (
	"errors"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"golang.org/x/crypto/bcrypt"
)

// Account is a user account that contains a password and a list of characters.
type Account struct {
	username   string
	ID         int
	Characters []*game.Character
	Password   string
}

// PasswordMatches returns if the given password matches the stored hashed password.
func (a *Account) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// HasCharacter returns of the account has the given character by name.
func (a *Account) HasCharacter(name string) bool {
	for _, c := range a.Characters {
		if c.Name == name {
			return true
		}
	}
	return false
}

// CreateCharacter creates a character with the name and archetype.
func (a *Account) CreateCharacter(name string, archetype id.UUID) (*game.Character, error) {
	if a.HasCharacter(name) {
		return nil, ErrCharacterExists
	}

	a.Characters = append(a.Characters, &game.Character{
		Name: name,
		WorldObject: game.WorldObject{
			ArchetypeID: archetype,
		},
	})

	return a.Characters[len(a.Characters)-1], nil
}

// DeleteCharacter deletes a given character by name.
func (a *Account) DeleteCharacter(name string) error {
	for i, ch := range a.Characters {
		if ch.Name == name {
			a.Characters = append(a.Characters[:i], a.Characters[i+1:]...)
			return nil
		}
	}
	return ErrCharacterDoesNotExist
}

var (
	ErrNotLoggedIn           = errors.New("not logged in")
	ErrBadPassword           = errors.New("bad password")
	ErrCharacterExists       = errors.New("character exists")
	ErrCharacterDoesNotExist = errors.New("character does not exist")
	ErrNoSuchArchetype       = errors.New("no such archetype")
)
