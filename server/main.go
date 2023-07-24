package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/kettek/morogue/game"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}

	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	// Load archetypes.
	var archetypes []game.Archetype
	{
		entries, err := os.ReadDir("archetypes")
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasSuffix(entry.Name(), ".json") {
				bytes, err := os.ReadFile(filepath.Join("archetypes", entry.Name()))
				if err != nil {
					log.Println(err)
					continue
				}
				var a game.Archetype
				err = json.Unmarshal(bytes, &a)
				if err != nil {
					log.Println(err)
					continue
				}
				archetypes = append(archetypes, a)
			}
		}
	}

	accounts, err := newAccounts("accounts")
	if err != nil {
		return err
	}

	u := newUniverse(accounts, archetypes)
	u.Run()

	ps := newSocketServer(u.clientChan, u.checkChan)

	s := &http.Server{
		Handler:      ps,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
