# morogue
morogue is a WIP multiplayer roguelike engine/game.

![WIP screenshot](screenshot.png)

## Developing
morogue uses the [gobl](https://github.com/kettek/gobl) build system. For development, you should start both the server and the client in watch mode:

```bash
go run . watch-server
```

```bash
go run . watch-client
```

Whenever watched files are changed, either program will be re-built and then re-executed.

To see the other available *gobl* commands, just issue:

```bash
go run .
```

## Architectural Notes

  * Networking uses websockets and all communication is done through [msgpack](https://msgpack.org/index.html) using go's unmarshal/marshal functionality.
    * Interfaces, such as Archetypes, Objects, Events, and Messages, employ the use of wrappers to safely marshal and unmarshal interfaces to and from their concrete types.
  * A **World** represents a contained game instance. It runs it its own goroutine and can have characters join or leave the world. Each individual "map" in a World is known as a **Location**.
  * Almost every distinct object in the world is of the **Object** type and contains a reference ID to an **Archetype** (and a cached pointer to said Archetype for efficiency). An Archetype contains the actual underlying data for an object, such as damage done, slots used, title. An Object is a "live" object that is used for actual world processing and interaction.
  * **Accounts** and their Characters are marshaled as JSON into a [bbolt](https://pkg.go.dev/go.etcd.io/bbolt#section-readme) database.
  * All Archetypes are defined as JSON files in various directories in the `archetypes` directory.
  * All Archetypes are defined and referenced by a UUIDv5 identifier. This identifier can be provided either by an ASCII string, an array of bytes, or by a human-readable string that is converted to the actual UUID. This human-readable string is written as `morogue:type:thing`, where *type* would be *armor*, *weapon*, *item*, *character*, or *tile*, and *thing* would be whatever the actual archetype is called.
  * Player controlled objects, such as Characters, receive commands from the player via a **Desire**. A desire can be to apply an item, drop an item, move in a direction, attack a target, and beyond. The result of a desire being processed will generally result in an **Event** being emitted to other clients or just the controlling player.
  * Events are generally used to represent something happening in a *Location*, such as a character equipping an item, something taking damage, and beyond.

## Adding New Archetypes Types
First add a new Archetype type, such as `FancyArchetype` in a corresponding file in the `game` subdir. Ensure it adheres to the `game.Archetype` interface and has a unique return value from Type(). Add a new Key, UUID, and its setup in init() within `id/namespaces.go`. Finally, add a new case entry to unmarshal in `*ArchetypeMessage UnmarshalMsgpack` within `net/message.go`.

Now you just need to update or add new features/logic wherever you wish the archetype to be used. At minimum, an object using the archetype can be picked up and dropped.