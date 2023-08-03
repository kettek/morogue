package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
)

// universe is the structure used to manage clients, accounts, character
// creation, and worlds. clients are moved to worlds for handling when
// they join and are no longer handled by the universe.
type universe struct {
	accounts               Accounts
	loggedInAccounts       []string
	clients                []*client
	clientChan             chan client
	clientRemoveChan       chan *client
	clientAddFromWorldChan chan *client
	checkChan              chan struct{}
	worlds                 []*world
	//
	data *Data
}

func newUniverse(accounts Accounts, data *Data) universe {
	return universe{
		accounts:               accounts,
		clientChan:             make(chan client, 10),
		checkChan:              make(chan struct{}, 10),
		clientRemoveChan:       make(chan *client, 10),
		clientAddFromWorldChan: make(chan *client, 10),
		data:                   data,
	}
}

// spinWorld adds the given world to the universe and starts the world's
// loop in a goroutine.
func (u *universe) spinWorld(w *world) {
	u.worlds = append(u.worlds, w)
	go w.loop(u.clientAddFromWorldChan, u.clientRemoveChan)
}

// Run starts the universe and returns a channel through which the world's
// processing can be stopped. The universe runs in a goroutine and handles
// client processing until a world takes over the client.
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
			case cl := <-u.clientRemoveChan:
				u.removeAccountLoggedIn(cl.account.username)
				if err := u.accounts.SaveAccount(cl.account); err != nil {
					log.Println(err)
				}
			case cl := <-u.clientAddFromWorldChan:
				cl.state = clientStateLoggedIn
				u.clients = append(u.clients, cl)
			}
		}
	}()
	return closeCh
}

// checkClients calls updateClient on all clients and removes
// any clients that are disconnected.
func (u *universe) checkClients() {
	i := 0
	for _, cl := range u.clients {
		if err := u.updateClient(cl); err == nil {
			u.clients[i] = cl
			i++
		} else if err != errRemoveClientFromUniverse {
			fmt.Println(err)
		}
	}
	for j := i; j < len(u.clients); j++ {
		u.clients[j] = nil
	}
	u.clients = u.clients[:i]
}

func (u *universe) checkName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	// TODO: Probably handle pottymouth names, if even possible.
	return nil
}

// loginClient sets a client to logged in and adds it to the
// logged in accounts slice.
func (u *universe) loginClient(cl *client) {
	cl.state = clientStateLoggedIn
	u.loggedInAccounts = append(u.loggedInAccounts, cl.account.username)
	// Send the available archetypes.
	cl.conn.Write(net.ArchetypesMessage{
		Archetypes: u.data.Archetypes,
	})
	// Send the player's characters.
	cl.conn.Write(net.CharactersMessage{
		Characters: cl.account.Characters,
	})
}

// updateClient processes all network messaging with a client,
// including logins, registrations, character creations,
// character deletions, world creation, and world joining.
func (u *universe) updateClient(cl *client) error {
	for {
		t := time.Now()
		select {
		case msg := <-cl.msgChan:
			switch m := msg.(type) {
			case net.LoginMessage:
				if cl.state != clientStateWaiting {
					cl.conn.Write(net.LoginMessage{
						ResultCode: 400,
					})
				} else {
					account, err := u.accounts.Account(m.User)
					if err != nil {
						if err == ErrNoUser {
							cl.conn.Write(net.LoginMessage{
								Result:     err.Error(),
								ResultCode: 404,
							})
						} else {
							cl.conn.Write(net.LoginMessage{
								Result:     err.Error(),
								ResultCode: 500,
							})
						}
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
						account, err := u.accounts.Account(m.User)
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
				if err := u.accounts.SaveAccount(cl.account); err != nil {
					log.Println(err)
				}
				cl.account = Account{}
				cl.state = clientStateWaiting
			case net.CreateCharacterMessage:
				if cl.state == clientStateSelectedCharacter {
					cl.conn.Write(net.JoinCharacterMessage{
						ResultCode: 400,
						Result:     ErrAlreadyJoined.Error(),
					})
				} else if cl.state != clientStateLoggedIn {
					cl.conn.Write(net.CreateCharacterMessage{
						ResultCode: 400,
						Result:     ErrNotLoggedIn.Error(),
					})
				} else {
					if !u.data.hasArchetype(m.Archetype) {
						cl.conn.Write(net.CreateCharacterMessage{
							ResultCode: 400,
							Result:     ErrNoSuchArchetype.Error(),
						})
					} else if err := u.checkName(m.Name); err != nil {
						cl.conn.Write(net.CreateCharacterMessage{
							ResultCode: 400,
							Result:     err.Error(),
						})
					} else if cl.account.HasCharacter(m.Name) {
						cl.conn.Write(net.CreateCharacterMessage{
							ResultCode: 400,
							Result:     ErrCharacterExists.Error(),
						})
					} else if char, err := cl.account.CreateCharacter(m.Name, m.Archetype); err != nil {
						cl.conn.Write(net.CreateCharacterMessage{
							ResultCode: 400,
							Result:     err.Error(),
						})
					} else {
						// Add the character's starting gear.
						if arch, ok := u.data.Archetype(char.Archetype).(game.CharacterArchetype); ok {
							for _, au := range arch.StartingObjects {
								if a := u.data.Archetype(au); a != nil {
									if o := game.CreateObjectFromArchetype(a); o != nil {
										// Use pickup to add the object to the character's inventory.
										char.Pickup(o)
										// Apply armor and weapons.
										if _, ok := o.(*game.Armor); ok {
											char.Apply(o)
										} else if _, ok := o.(*game.Weapon); ok {
											char.Apply(o)
										}
									}
								}
							}
						}
						// Save the changes.
						u.accounts.SaveAccount(cl.account)
						// Let 'em know it went ok
						cl.conn.Write(net.CreateCharacterMessage{
							ResultCode: 200,
							Result:     fmt.Sprintf("%s takes form", m.Name),
						})
						// Re-send the player's characters.
						cl.conn.Write(net.CharactersMessage{
							Characters: cl.account.Characters,
						})
					}
				}
			case net.DeleteCharacterMessage:
				if cl.state == clientStateSelectedCharacter {
					cl.conn.Write(net.JoinCharacterMessage{
						ResultCode: 400,
						Result:     ErrAlreadyJoined.Error(),
					})
				} else if cl.state != clientStateLoggedIn {
					cl.conn.Write(net.DeleteCharacterMessage{
						ResultCode: 400,
						Result:     ErrNotLoggedIn.Error(),
					})
				} else {
					if err := cl.account.DeleteCharacter(m.Name); err != nil {
						cl.conn.Write(net.DeleteCharacterMessage{
							ResultCode: 400,
							Result:     err.Error(),
						})
					} else {
						u.accounts.SaveAccount(cl.account)
						// Let 'em know it went ok
						cl.conn.Write(net.DeleteCharacterMessage{
							ResultCode: 200,
							Result:     fmt.Sprintf("%s is no more", m.Name),
						})
						// Re-send the player's characters.
						cl.conn.Write(net.CharactersMessage{
							Characters: cl.account.Characters,
						})
					}
				}
			case net.JoinCharacterMessage:
				if cl.state == clientStateSelectedCharacter {
					cl.conn.Write(net.JoinCharacterMessage{
						ResultCode: 400,
						Result:     ErrAlreadyJoined.Error(),
					})
				} else if cl.state != clientStateLoggedIn {
					cl.conn.Write(net.JoinCharacterMessage{
						ResultCode: 400,
						Result:     ErrNotLoggedIn.Error(),
					})
				} else {
					if !cl.account.HasCharacter(m.Name) {
						cl.conn.Write(net.JoinCharacterMessage{
							ResultCode: 400,
							Result:     ErrCharacterDoesNotExist.Error(),
						})
					} else {
						// Mark character as desired for client.
						cl.character = m.Name
						cl.state = clientStateSelectedCharacter
						// Let the character know we've considered them as joined for that character.
						cl.conn.Write(net.JoinCharacterMessage{
							ResultCode: 200,
						})
						// And send the current list of worlds.
						if t.Sub(cl.lastWorldsSent) <= 2*time.Second {
							cl.conn.Write(net.WorldsMessage{
								ResultCode: 429,
								Result:     ErrTooSoon.Error(),
							})
						} else {
							cl.conn.Write(net.WorldsMessage{
								ResultCode: 200,
								Worlds:     u.getWorldsInfos(),
							})
							cl.lastWorldsSent = t
						}
					}
				}
			case net.UnjoinCharacterMessage:
				if cl.state != clientStateSelectedCharacter {
					cl.conn.Write(net.UnjoinCharacterMessage{
						ResultCode: 400,
						Result:     ErrNotJoined.Error(),
					})
				} else {
					cl.character = ""
					cl.state = clientStateLoggedIn
					cl.conn.Write(net.UnjoinCharacterMessage{
						ResultCode: 200,
					})
				}
			case net.WorldsMessage:
				if cl.state < clientStateSelectedCharacter {
					cl.conn.Write(net.WorldsMessage{
						ResultCode: 400,
						Result:     ErrWrongState.Error(),
					})
				} else if t.Sub(cl.lastWorldsSent) <= 2*time.Second {
					cl.conn.Write(net.WorldsMessage{
						ResultCode: 429,
						Result:     ErrTooSoon.Error(),
					})
				} else {
					cl.conn.Write(net.WorldsMessage{
						ResultCode: 200,
						Worlds:     u.getWorldsInfos(),
					})
					cl.lastWorldsSent = t
				}
			case net.CreateWorldMessage:
				if cl.state < clientStateSelectedCharacter || cl.state > clientStateSelectedCharacter {
					cl.conn.Write(net.CreateWorldMessage{
						ResultCode: 400,
						Result:     ErrWrongState.Error(),
					})
				} else {
					// TODO: Throttle this as well.
					w := newWorld(u.data)
					if m.Password != "" {
						w.info.Private = true
						w.password = m.Password
					}
					w.info.Name = m.Name
					u.spinWorld(w)
					cl.conn.Write(net.JoinWorldMessage{
						ResultCode: 200,
						World:      w.info.ID,
					})
					w.clientChan <- cl
					return errRemoveClientFromUniverse
				}
			case net.JoinWorldMessage:
				if cl.state < clientStateSelectedCharacter || cl.state > clientStateSelectedCharacter {
					cl.conn.Write(net.JoinWorldMessage{
						ResultCode: 400,
						Result:     ErrWrongState.Error(),
					})
				} else if w, err := u.getWorld(m.World); err != nil {
					cl.conn.Write(net.JoinWorldMessage{
						ResultCode: 400,
						Result:     err.Error(),
					})
				} else if w.password != m.Password {
					cl.conn.Write(net.JoinWorldMessage{
						ResultCode: 400,
						Result:     ErrBadPassword.Error(),
					})
				} else {
					cl.conn.Write(net.JoinWorldMessage{
						ResultCode: 200,
					})
					w.clientChan <- cl
					return errRemoveClientFromUniverse
				}
			}
		case err := <-cl.closedChan:
			if cl.account.username != "" {
				u.removeAccountLoggedIn(cl.account.username)
				if err := u.accounts.SaveAccount(cl.account); err != nil {
					log.Println(err)
				}
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

// getWorldsInfos returns the infos of all worlds.
func (u *universe) getWorldsInfos() (worlds []game.WorldInfo) {
	for _, w := range u.worlds {
		worlds = append(worlds, w.info)
	}
	return
}

func (u *universe) getWorld(uuid id.UUID) (*world, error) {
	for _, w := range u.worlds {
		if w.info.ID == uuid {
			return w, nil
		}
	}
	return nil, ErrWorldDoesNotExist
}

var (
	ErrWrongState               = errors.New("message sent in wrong state")
	ErrTooSoon                  = errors.New("message sent too soon, please wait and try again")
	ErrNotJoined                = errors.New("character not joined")
	ErrAlreadyJoined            = errors.New("character is already joined")
	ErrUserLoggedIn             = errors.New("user is logged in")
	ErrWorldDoesNotExist        = errors.New("world does not exist")
	errRemoveClientFromUniverse = errors.New("this is not an error lol")
)
