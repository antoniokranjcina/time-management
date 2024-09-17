package domain

type LocationRepository interface {
	Save(location *Location) (*Location, error)
	GetAll() ([]Location, error)
	GetById(id string) (*Location, error)
	Update(id string, name string) (*Location, error)
	Delete(id string) error
}
