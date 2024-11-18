package server

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// Accounts is an interface for loading, creating, and saving accounts.
type Accounts interface {
	Account(username string) (account Account, err error)
	NewAccount(username string, password string) error
	SaveAccount(account Account) error
	Buckets() (buckets []string)
	ListBucket(bucket string) (list []string)
	DumpBytes(bucket, data string) []byte
}

type accounts struct {
	db *bolt.DB
}

// NewAccounts creates a new accounts database with the given path.
func NewAccounts(path string) (Accounts, error) {
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

// Account returns an account with the given username.
func (a *accounts) Account(username string) (account Account, err error) {
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

// NewAccount creates a new account.
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
	// remove unneeded stuff from account
	// FIXME: Move this to universe
	for _, ch := range account.Characters {
		ch.Desire = nil
		ch.LastDesire = nil
	}

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

// Buckets returns the bolt buckets of the accounts.
func (a *accounts) Buckets() (buckets []string) {
	a.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
		return nil
	})
	return
}

// ListBucket returns the a list of buckets.
func (a *accounts) ListBucket(bucket string) (list []string) {
	a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.ForEach(func(k, v []byte) error {
			list = append(list, string(k))
			return nil
		})
		return nil
	})
	return
}

// DumpBytes dumps some bytes.
func (a *accounts) DumpBytes(bucket, data string) (d []byte) {
	a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		d = b.Get([]byte(data))
		return nil
	})
	return
}

// Accounts-related errors.
var (
	ErrNoUser     = errors.New("no such user")
	ErrUserExists = errors.New("user exists")
)
