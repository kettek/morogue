package main

import (
	"fmt"

	"github.com/kettek/morogue/net"
)

type universe struct {
	accounts   Accounts
	clients    []client
	clientChan chan client
	checkChan  chan struct{}
	worlds     []world
}

func newUniverse(accounts Accounts) universe {
	return universe{
		accounts:   accounts,
		clientChan: make(chan client, 10),
		checkChan:  make(chan struct{}, 10),
	}
}

func (u *universe) spinWorld() *world {
	w := &world{
		quitChan: make(chan struct{}),
	}

	return w
}

func (u *universe) Run() chan struct{} {
	closeCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-closeCh:
				// TODO: Cleanup worlds.
				return
			case client := <-u.clientChan:
				u.clients = append(u.clients, client)
			case <-u.checkChan:
				u.checkClients()
			}
		}
	}()
	return closeCh
}

func (u *universe) checkClients() {
	i := 0
	for _, cl := range u.clients {
		if err := u.updateClient(cl); err == nil {
			u.clients[i] = cl
			i++
		} else {
			fmt.Println(err)
		}
	}
	for j := i; j < len(u.clients); j++ {
		u.clients[j] = client{}
	}
	u.clients = u.clients[:i]
}

func (u *universe) updateClient(cl client) error {
	fmt.Println("updateClient")
	for {
		select {
		case msg := <-cl.msgChan:
			switch m := msg.(type) {
			case net.LoginMessage:
				if cl.state != clientStateWaiting {
					cl.conn.Write(net.LoginMessage{
						ResultCode: 400,
					})
				} else {
					account, err := u.accounts.GetAccount(m.User)
					if err != nil {
						cl.conn.Write(net.LoginMessage{
							Result:     err.Error(),
							ResultCode: 404,
						})
					} else {
						if !account.PasswordMatches(m.Password) {
							cl.conn.Write(net.LoginMessage{
								Result:     ErrBadPassword.Error(),
								ResultCode: 403,
							})
						} else {
							cl.conn.Write(net.LoginMessage{
								ResultCode: 200,
							})
							cl.account = account
							cl.state = clientStateLoggedIn
						}
					}
				}
			case net.RegisterMessage:
				if cl.state != clientStateWaiting {
					cl.conn.Write(net.RegisterMessage{
						ResultCode: 400,
					})
				} else {
					err := u.accounts.NewAccount(m.User, m.Password)
					if err != nil {
						cl.conn.Write(net.RegisterMessage{
							Result:     err.Error(),
							ResultCode: 404,
						})
					} else {
						account, err := u.accounts.GetAccount(m.User)
						if err != nil {
							cl.conn.Write(net.RegisterMessage{
								Result:     err.Error(),
								ResultCode: 404,
							})
						} else {
							cl.conn.Write(net.RegisterMessage{
								ResultCode: 200,
							})
							cl.account = account
							cl.state = clientStateLoggedIn
						}
					}
				}
			}
		case err := <-cl.closedChan:
			fmt.Println("client yeeted", err)
			return err
		default:
			return nil
		}
	}
}
