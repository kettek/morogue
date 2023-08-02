package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	data := &Data{}
	if err := data.loadArchetypes(); err != nil {
		return err
	}

	accounts, err := newAccounts("accounts")
	if err != nil {
		return err
	}

	u := newUniverse(accounts, data)
	u.Run()

	ps := newSocketServer(u.clientChan, u.checkChan)

	// Allow access to archetypes via archetypes subdir.
	ps.serveMux.Handle("/archetypes/", http.StripPrefix("/archetypes/", http.FileServer(http.Dir("./archetypes"))))

	ps.serveMux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

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
