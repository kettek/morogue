package game

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type Inventory struct {
	inventory *[]game.Object
	container *widget.Container
}

func (inv *Inventory) Init(container *widget.Container) {
	inv.container = container
}

func (inv *Inventory) SetInventory(inventory *[]game.Object) {
	inv.inventory = inventory
}

func (inv *Inventory) Refresh() {
	// TODO: Refresh inventory
}

func (inv *Inventory) Draw(ctx ifs.DrawContext) {
	// TODO: Draw
}
