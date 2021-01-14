package job

type Publisher interface {
	Publish(job *Job) error
}
