package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var app string
var data any

func Init[K any](a string, v K) K {
	app = a
	if err := Load(v); err != nil {
		if err := Save(); err != nil {
			panic(err)
		}
	}
	data = v
	return v
}

func Save() error {
	return save(data)
}

func save[K any](v K) error {
	d, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	p := filepath.Join(d, app)

	if err := os.MkdirAll(p, 0755); err != nil {
		return err
	}

	c := filepath.Join(p, "config.json")

	f, err := os.Create(c)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(v)
}

func Load[K any](v K) error {
	d, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	p := filepath.Join(d, app)

	c := filepath.Join(p, "config.json")

	f, err := os.Open(c)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}
