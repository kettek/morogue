package ifs

type State interface {
	Leave() error
	Return(interface{}) error
	Begin() error
	Update(RunContext) error
	Draw(DrawContext)
	End() (interface{}, error)
}
