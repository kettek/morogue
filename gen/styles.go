package gen

type MapConfig struct {
	Width  int
	Height int
}

type Styler struct {
	Passes      int
	ProcessPass func(cfg Config, pass int) error
}

type Config struct {
	Width, Height int
	Cell          func(x, y int) Cell
	SetCell       func(x, y int, cell Cell)
}

var Styles = map[string]Styler{}

func RegisterStyle(name string, cfg Styler) {
	Styles[name] = cfg
}
