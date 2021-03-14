package container

// Provider is used for providing services and parameters.
type Provider = func() (interface{}, error)

// Decorator decorates service passed as a second argument.
// The first argument is the ID of a given service.
type Decorator = func(string, interface{}) (interface{}, error)
