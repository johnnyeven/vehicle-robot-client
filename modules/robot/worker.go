package robot

type Worker interface {
	WorkerID() string
	Start()
	Restart() error
	Stop() error
}
