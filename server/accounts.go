package main

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

type Accounts interface {
	GetAccount(username string) (account Account, err error)
	NewAccount(username string, password string) error
	SaveAccount(account Account) error
}

type accounts struct {
	db *bolt.DB
}

func newAccounts(path string) (*accounts, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		return err
	})

	a := &accounts{
		db: db,
	}

	return a, err
}

func (a *accounts) GetAccount(username string) (account Account, err error) {
	err = a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		data := b.Get([]byte(username))

		if data == nil {
			return ErrNoUser
		}

		err := json.Unmarshal(data, &account)

		return err
	})

	return
}

func (a *accounts) NewAccount(username string, password string) error {
	err := a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		data := b.Get([]byte(username))
		if data != nil {
			return ErrUserExists
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			return err
		}

		id, _ := b.NextSequence()

		account := Account{
			Password: string(bytes),
			ID:       int(id),
		}

		buf, err := json.Marshal(&account)
		if err != nil {
			return err
		}

		return b.Put([]byte(username), buf)
	})

	return err
}

// SaveAccount saves the given Account. This must be an account acquired from GetAccount.
func (a *accounts) SaveAccount(account Account) error {
	err := a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))

		buf, err := json.Marshal(&account)
		if err != nil {
			return err
		}

		return b.Put([]byte(account.username), buf)
	})
	return err
}

var (
	ErrNoUser     = errors.New("no such user")
	ErrUserExists = errors.New("user exists")
)
