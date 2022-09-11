package main

type Photo struct {
	ID     int
	UserID int
	Path   string
}

// attempt to abstract storage
type StMem struct {
	items []*Photo
}

func NewStorage() *StMem {
	return &StMem{
		items: make([]*Photo, 0, 10),
	}
}

func (st *StMem) Add(p *Photo) error {
	st.items = append(st.items, p)
	return nil
}

func (st *StMem) GetPhotos(userID int) ([]*Photo, error) {
	return st.items, nil
}
