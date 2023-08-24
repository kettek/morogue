package gen

import "math/rand"

type MinMax [2]int

func (m MinMax) Min() int {
	return m[0]
}

func (m MinMax) Max() int {
	return m[1]
}

func (m MinMax) Roll() int {
	if m[0] == m[1] {
		return m[0]
	}
	return m[0] + rand.Intn(m[1]-m[0])
}
