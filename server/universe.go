package main

import (
	"errors"
	"fmt"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
)

type universe struct {
	accounts         Accounts
	loggedInAccounts []string
	clients          []*client
	clientChan       chan client
	checkChan        chan struct{}
	worlds           []world
	//
	archetypes []game.Archetype
}

func newUniverse(accounts Accounts, archetypes []game.Archetype) universe {
	return universe{
		accounts:   accounts,
		clientChan: make(chan client, 10),
		checkChan:  make(chan struct{}, 10),
		archetypes: archetypes,
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
				u.clients = append(u.clients, &client)
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
		u.clients[j] = nil
	}
	u.clients = u.clients[:i]
}

func (u *universe) hasArchetype(uuid id.UUID) bool {
	for _, a := range u.archetypes {
		if a.UUID == uuid {
			return true
		}
	}
	return false
}

func (u *universe) checkName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	// TODO: Probably handle pottymouth names, if even possible.
	return nil
}

func (u *universe) loginClient(cl *client) {
	cl.state = clientStateLoggedIn
	u.loggedInAccounts = append(u.loggedInAccounts, cl.account.username)
	// Send the available archetypes.
	cl.conn.Write(net.ArchetypesMessage{
		Archetypes: u.archetypes,
	})
	// Send the player's characters.
	cl.conn.Write(net.CharactersMessage{
		Characters: cl.account.Characters,
	})
}

func (u *universe) updateClient(cl *client) error {
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
						} else if u.isAccountLoggedIn(m.User) {
							cl.conn.Write(net.LoginMessage{
								Result:     ErrUserLoggedIn.Error(),
								ResultCode: 403,
							})
						} else {
							cl.conn.Write(net.LoginMessage{
								ResultCode: 200,
							})
							account.username = m.User
							cl.account = account
							u.loginClient(cl)
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
							account.username = m.User
							cl.account = account
							u.loginClient(cl)
						}
					}
				}
			case net.LogoutMessage:
				u.removeAccountLoggedIn(cl.account.username)
				cl.account = Account{}
				cl.state = clientStateWaiting
			case net.CreateCharacterMessage:
				if cl.state != clientStateLoggedIn {
					cl.conn.Write(net.RegisterMessage{
						ResultCode: 400,
						Result:     "not logged in",
					})
				} else {
					if !u.hasArchetype(m.Archetype) {
						cl.conn.Write(net.RegisterMessage{
							ResultCode: 400,
							Result:     "no such archetype",
						})
					} else if err := u.checkName(m.Name); err != nil {
						cl.conn.Write(net.RegisterMessage{
							ResultCode: 400,
							Result:     err.Error(),
						})
					} else if cl.account.HasCharacter(m.Name) {
						cl.conn.Write(net.RegisterMessage{
							ResultCode: 400,
							Result:     "character with given name already exists",
						})
					} else {
						fmt.Println("Handle character creation", m)
					}
				}
			}
		case err := <-cl.closedChan:
			if cl.account.username != "" {
				u.removeAccountLoggedIn(cl.account.username)
			}
			fmt.Println("client yeeted", err)
			return err
		default:
			return nil
		}
	}
}

func (u *universe) isAccountLoggedIn(username string) bool {
	for _, user := range u.loggedInAccounts {
		if user == username {
			return true
		}
	}
	return false
}

func (u *universe) removeAccountLoggedIn(username string) {
	fmt.Println("removing", username, u.loggedInAccounts)
	for i, user := range u.loggedInAccounts {
		if user == username {
			u.loggedInAccounts = append(u.loggedInAccounts[:i], u.loggedInAccounts[i+1:]...)
			return
		}
	}
}

var (
	ErrUserLoggedIn = errors.New("user is logged in")
)
