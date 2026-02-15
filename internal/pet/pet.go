package pet

// Pet is the domain model for a pet.
type Pet struct {
	ID   int64
	Name string
	Tag  *string
}
