package gen

const StyleBox string = "default-box"

func init() {
	styler := Styler{
		Passes: 1,
		ProcessPass: func(cfg Config, pass int) error {
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
