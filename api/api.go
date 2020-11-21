package api

// An Action is a unit of execution for the API
type Action interface {
	Execute() error
}
