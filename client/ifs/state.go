package ifs

type State interface {
	Leave() error
	Return(interface{}) error
	Begin(RunContext) error
	Update(RunContext) error
	Draw(DrawContext)
	End() (interface{}, error)
}
