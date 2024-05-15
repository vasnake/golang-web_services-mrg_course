package storage

type Photo struct {
	ID     int
	UserID int
	Path   string
}

// not interface, but we getting there

type StMem struct {
	items []*Photo
}

// factory, init storage
func NewStorage() *StMem {
	return &StMem{
		items: make([]*Photo, 0, 16),
	}
}

func (st *StMem) Add(p *Photo) error {
	st.items = append(st.items, p)
	return nil
}

func (st *StMem) GetPhotos(userID int) ([]*Photo, error) {
	return st.items, nil
}
