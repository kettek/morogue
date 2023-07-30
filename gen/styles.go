package gen

type MapConfig struct {
	Width  int
	Height int
}

type Styler struct {
	Passes      int
	ProcessPass func(cfg Config, pass int) error
}

var Styles = map[string]Styler{}

func RegisterStyle(name string, cfg Styler) {
	Styles[name] = cfg
}
