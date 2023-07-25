package main

import (
	"fmt"
	"os"

	"github.com/kettek/morogue/id"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Syntax: ./uuidgen [archetype|item] <name of item>")
		return
	}

	var namespace id.UUID

	switch os.Args[1] {
	case "archetype":
		namespace = id.Archetype
	case "mob":
		namespace = id.Mob
	case "item":
		namespace = id.Item
	default:
		fmt.Println("namespace must be one of: archetype, mob, item")
		return
	}

	uuid, err := id.UID(namespace, os.Args[2])
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%s:%s", id.NamespaceToKey[namespace], os.Args[2]), "=", uuid)
}
