package gen

const StyleBox string = "default-box"

type ConfigBox struct {
	Width   int
	Height  int
	Cell    func(x, y int) Cell
	SetCell func(x, y int, cell Cell)
}

func init() {
	styler := Styler{
		Passes: 1,
		ProcessPass: func(cf Config, pass int) error {
			cfg := cf.(ConfigBox)
			for x := 0; x < cfg.Width; x++ {
				for y := 0; y < cfg.Height; y++ {
					if x == 0 || y == 0 || x == cfg.Width-1 || y == cfg.Height-1 {
						c := cfg.Cell(x, y)
						if !c.Flags().Has("blocked") {
							c.SetFlags([]string{"blocked"})
							cfg.SetCell(x, y, c)
						}
					}
				}
			}
			return nil
		},
	}
	RegisterStyle(StyleBox, styler)
}
