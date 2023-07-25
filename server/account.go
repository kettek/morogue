package main

import (
	"errors"

	"github.com/kettek/morogue/game"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	username   string
	ID         int
	Characters []game.Character
	Password   string
}

func (a *Account) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

func (a *Account) HasCharacter(name string) bool {
	for _, c := range a.Characters {
		if c.Name == name {
			return true
		}
	}
	return false
}

var (
	ErrBadPassword = errors.New("bad password")
)
