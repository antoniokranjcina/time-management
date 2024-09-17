package domain

type Location struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt uint64 `json:"created_at"`
}

// NewLocation Factory method to create a Location
func NewLocation(id, name string, createdAt uint64) *Location {
	return &Location{
		Id:        id,
		Name:      name,
		CreatedAt: createdAt,
	}
}
