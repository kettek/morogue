package gen

func Generate(styler Styler, cfg Config) error {
	for i := 0; i < styler.Passes; i++ {
		err := styler.ProcessPass(cfg, i)
		if err != nil {
			return err
		}
	}
	return nil
}
