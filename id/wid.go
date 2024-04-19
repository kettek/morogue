package id

// WID is UUIDv4 that represents world-based things.
type WID int

// WIDGenerator generates WIDs.
type WIDGenerator struct {
	top WID
}

// Next returns the next WID (top+1) and increases the internal top WID.
func (w *WIDGenerator) Next() WID {
	w.top++
	return w.top
}

// Top returns the "top" WID, which corresponds to the last WID generated.
func (w *WIDGenerator) Top() WID {
	return w.top
}
