package gen

func init() {
	RegisterStyle("default-void", Styler{
		Passes: 1,
		ProcessPass: func(cfg Config, pass int) error {
			return nil
		},
	})
}
