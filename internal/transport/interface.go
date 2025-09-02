package transport

type Server interface {
	Start(address string) error
	Stop() error
}

type TransportData interface {
	Path() string
	Method() string
}
