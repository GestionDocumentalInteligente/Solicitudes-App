package defs

type Server interface {
	Run() error
	SetRouter(router interface{}) error
}

type Config interface {
	GetServerName() string
	GetServerHost() string
	GetServerPort() int
	GetServerID() string
	GetServerAddress() string
	GetConsulAddress() string
	GetRouter() interface{}
	Validate() error
}
