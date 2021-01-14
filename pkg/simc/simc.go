package simc

type Configuration string

type Result struct {
	Data map[string]interface{}
}

type Simulator interface {
	Simulate(Configuration) (*Result, error)
}
