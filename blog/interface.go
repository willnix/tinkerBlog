package blog

type Blogger interface {
	FindById(id string) (*Entry, error)
	LatestEntries() ([]*Entry, error)
	Delete(id string) error
	Save(*Entry) error
}
