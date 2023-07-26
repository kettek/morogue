package gen

type coord [2]int

type Flags []string

func (f Flags) Has(v string) bool {
	for _, fl := range f {
		if fl == v {
			return true
		}
	}
	return false
}

type Cell interface {
	Value() int
	SetValue(v int)
	Flags() Flags
	SetFlags(v Flags)
	Data() interface{}
	SetData(d interface{})
}

type cell struct {
	value int
	flags []string
	data  interface{}
}

func (c *cell) Value() int {
	return c.value
}
func (c *cell) SetValue(v int) {
	c.value = v
}
func (c *cell) Flags() []string {
	return c.flags
}
func (c *cell) SetFlags(f []string) {
	c.flags = f
}
func (c *cell) Data() interface{} {
	return c.data
}
func (c *cell) SetData(d interface{}) {
	c.data = d
}
