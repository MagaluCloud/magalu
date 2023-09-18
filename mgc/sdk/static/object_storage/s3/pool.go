package s3

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var poolLoggerInstance *zap.SugaredLogger

func poolLogger() *zap.SugaredLogger {
	if poolLoggerInstance == nil {
		poolLoggerInstance = logger().Named("delete")
	}
	return poolLoggerInstance
}

type ThreadPool struct {
	maxWorkers int
	queueSize  int64
	waitGroup  sync.WaitGroup
	// A pool of workers channels that are registered with the dispatcher
	workerPool chan chan Job
	// A buffered channel that we can send work requests on.
	jobQueue chan Job
	// Channel used to stop all the workers
	close chan bool
}

func NewThreadPool(maxWorkers int, queueSize int64) *ThreadPool {
	tp := &ThreadPool{
		maxWorkers: maxWorkers,
		queueSize:  queueSize,
		waitGroup:  sync.WaitGroup{},
		workerPool: make(chan chan Job, maxWorkers),
		jobQueue:   make(chan Job, queueSize),
	}
	tp.Setup()
	return tp
}

// Creates the workers and start listening on the jobQueue
func (t *ThreadPool) Setup() {
	// starting n number of workers
	for i := 0; i < t.maxWorkers; i++ {
		worker := NewWorker(t.workerPool, t.waitGroup.Done, t.close)
		worker.Start()
	}

	go t.dispatch()
}

func (t *ThreadPool) Wait() {
	t.waitGroup.Wait()
}

// Enqueue submits the job to available worker
func (t *ThreadPool) Enqueue(job Job) error {
	// Add the task to the job queue
	if len(t.jobQueue) == int(t.queueSize) {
		return fmt.Errorf("Queue is full: %d scheduled jobs\n", t.queueSize)
	}
	t.jobQueue <- job
	t.waitGroup.Add(1)
	return nil
}

func (t *ThreadPool) dispatch() {
	for {
		select {
		case job := <-t.jobQueue:
			// a job request has been received
			func(job Job) {
				//Find a worker for the job
				jobChannel := <-t.workerPool
				//Submit job to the worker
				jobChannel <- job
			}(job)
		case <-t.close:
			// Close thread threadpool
			return
		}
	}
}

// Job represents the job to be run
type Job struct {
	Execute func() error
}

// Worker represents the worker that executes the job
type Worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	done       func()
	quit       chan bool
}

func NewWorker(workerPool chan chan Job, done func(), quit chan bool) Worker {
	return Worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		done:       done,
		quit:       quit,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				if err := job.Execute(); err != nil {
					poolLogger().Errorf("Error executing S3 job: %w", err)
				}
				w.done()
			case <-w.quit:
				w.done()
				return
			}
		}
	}()
}
