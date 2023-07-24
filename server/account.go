package main

import (
	"errors"

	"github.com/kettek/morogue/game"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID         int
	Characters []game.Character
	Password   string
}

func (a *Account) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

var (
	ErrBadPassword = errors.New("bad password")
)
