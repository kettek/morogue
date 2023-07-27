package id

// WID is UUIDv4 that represents world-based things.
type WID int

type WIDGenerator struct {
	top WID
}

func (w *WIDGenerator) Next() WID {
	w.top++
	return w.top
}
