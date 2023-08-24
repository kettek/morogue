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

	"github.com/kettek/morogue/server"
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

	data := &server.Data{}
	if err := data.LoadArchetypes(); err != nil {
		return err
	}
	log.Println(len(data.Archetypes), "archetypes")
	if err := data.LoadPlaces(); err != nil {
		return err
	}
	log.Println(len(data.Places), "places")
	if err := data.LoadFixtures(); err != nil {
		return err
	}
	log.Println(len(data.Fixtures), "fixtures")

	accounts, err := server.NewAccounts("accounts")
	if err != nil {
		return err
	}

	u, clientChan, checkChan := server.NewUniverse(accounts, data)
	u.Run()

	ps := server.NewSocketServer(clientChan, checkChan)

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
