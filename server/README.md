# Morogue Server
The morogue server supplies an HTTP and WebSockets server. The websockets server interfaces with a running morogue universe.

A universe contains:
  * accounts
    * characters
  * worlds
    * locations

The universe handles creating accounts and worlds. If a client has selected a character and either creates or joins a world, the client is sent to the world and thereafter the world handles all network communication. Once the world is done with the client, it sends it back to the universe to be managed by it or to be removed by it.

An account is a user-created account. Each account can store multiple characters. These characters are saved between worlds unless a permadeath server is chosen.

A world is a goroutine that contains locations and clients. Each location contains characters, mobs, objects, and, of course, the cells that make up a location. In more traditional terms, a world is a self-contained game state, and a location is a map or level.
