package job

type Consumer interface {
	Consume() (*Job, error)
}
