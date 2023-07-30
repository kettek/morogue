package gen

import (
	"fmt"
	"math/rand"
)

const StyleRooms string = "default-rooms"

type ConfigRooms struct {
	Width           int
	Height          int
	MinRoomSize     int
	MaxRoomSize     int
	MaxRooms        int
	JoinSharedWalls bool
	OverlapPadding  int
	Cell            func(x, y int) Cell
	SetCell         func(x, y int, cell Cell)
}

func isAreaOpen(cell func(x, y int) Cell, x1, y1, x2, y2 int) bool {
	var minx, maxx int
	var miny, maxy int

	if x1 < x2 {
		minx = x1
		maxx = x2
	} else if x1 > x2 {
		minx = x2
		maxx = x1
	}
	if y1 < y2 {
		miny = y1
		maxy = y2
	} else if y1 > y2 {
		miny = y2
		maxy = y1
	}

	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			if c := cell(x, y); c == nil || c.Flags().Has("room-floor") {
				return false
			}
		}
	}
	return true
}

func FillRange(getCell func(x, y int) Cell, setCell func(x, y int, cell Cell), x1, y1, x2, y2 int, flag ...string) {
	var minx, maxx int
	var miny, maxy int

	if x1 < x2 {
		minx = x1
		maxx = x2
	} else if x1 > x2 {
		minx = x2
		maxx = x1
	}
	if y1 < y2 {
		miny = y1
		maxy = y2
	} else if y1 > y2 {
		miny = y2
		maxy = y1
	}

	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			if c := getCell(x, y); c != nil {
				c.SetFlags(append(c.Flags(), flag...))
				setCell(x, y, c)
			}
		}
	}
}

func init() {
	styler := Styler{
		Passes: 1,
		ProcessPass: func(cf Config, pass int) error {
			cfg, ok := cf.(ConfigRooms)
			if !ok {
				return ErrWrongConfig
			}

			tries := 0
			success := 0
			for done := false; !done && tries < 20; tries++ {
				size := int(int32(cfg.MinRoomSize) + rand.Int31n(int32(cfg.MaxRoomSize)))
				x := int(rand.Int31n(int32(cfg.Width)))
				y := int(rand.Int31n(int32(cfg.Height)))

				if !isAreaOpen(cfg.Cell, x, y, x, y) {
					continue
				}

				x1 := x - size/2
				x2 := x + size/2
				y1 := y - size/2
				y2 := y + size/2

				hasTL := false
				hasTR := false
				hasBL := false
				hasBR := false
				sides := 0
				// If there are any free areas from x&y +/- size/2, then we can consider it valid.
				if isAreaOpen(cfg.Cell, x+cfg.OverlapPadding, y+cfg.OverlapPadding, x1-1, y1-1) {
					hasTL = true
					sides++
				}
				if isAreaOpen(cfg.Cell, x+cfg.OverlapPadding, y+cfg.OverlapPadding, x2+1, y1-1) {
					hasTR = true
					sides++
				}
				if isAreaOpen(cfg.Cell, x-cfg.OverlapPadding, y-cfg.OverlapPadding, x1-1, y2+1) {
					hasBL = true
					sides++
				}
				if isAreaOpen(cfg.Cell, x-cfg.OverlapPadding, y-cfg.OverlapPadding, x2+1, y2+1) {
					hasBR = true
					sides++
				}

				if sides <= 1 {
					continue
				}

				c := cfg.Cell(x, y)
				c.SetFlags(append(c.Flags(), "room-wall"))
				cfg.SetCell(x, y, c)

				if hasTL {
					FillRange(cfg.Cell, cfg.SetCell, x, y, x1, y1, "room-floor", "room", fmt.Sprintf("room#%d", success))
				}
				if hasTR {
					FillRange(cfg.Cell, cfg.SetCell, x, y, x2, y1, "room-floor", "room", fmt.Sprintf("room#%d", success))
				}
				if hasBL {
					FillRange(cfg.Cell, cfg.SetCell, x, y, x1, y2, "room-floor", "room", fmt.Sprintf("room#%d", success))
				}
				if hasBR {
					FillRange(cfg.Cell, cfg.SetCell, x, y, x2, y2, "room-floor", "room", fmt.Sprintf("room#%d", success))
				}

				success++
				tries = 0
				fmt.Println(success, tries)
				if success >= cfg.MaxRooms {
					break
				}
			}

			// Join any cells that have more than 1 wall.
			if cfg.JoinSharedWalls {
				for x := 0; x < cfg.Width-1; x++ {
					for y := 0; y < cfg.Height-1; y++ {
						if c := cfg.Cell(x, y); c != nil {
							count := 0
							flags := c.Flags()
							for _, s := range flags {
								if s == "room-wall" {
									count++
								}
							}
							if count > 1 {
								i := 0
								for _, s := range flags {
									if s != "room-wall" {
										flags[i] = s
										i++
									}
								}
								flags = flags[:i]
								c.SetFlags(flags)
								cfg.SetCell(x, y, c)
							}
						}
					}
				}
			}

			return nil
		},
	}
	RegisterStyle(StyleRooms, styler)
}
