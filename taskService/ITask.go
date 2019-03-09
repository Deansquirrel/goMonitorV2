package taskService

type ITask interface {
	StartTask() error
	StartJob(id string) error
	StopJob(id string) error
	RefreshConfig() error
}
