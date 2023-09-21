package core

// Interface for all messages that could be sent to thread pool
type TPMessage interface {
	// This is necessary to differentiate between exiting (true) and non-exiting (false) messages
	execute() bool
}

// Job type to run in thread pool
type TPJob func()

// Run TPJob as normal function, but mark as non-exiting message
func (j TPJob) execute() bool {
	j()
	return false
}

var _ TPMessage = (TPJob)(nil)

type exitMessage func()

// Mark message as exiting message
func (em exitMessage) execute() bool {
	em()
	return true
}

// Allocates Threads to solve messages in Queue
type ThreadPool struct {
	finished     bool
	workerN      int
	messageQueue chan TPMessage
}

func NewThreadPool(workerN int, queueSize int64) *ThreadPool {
	tp := &ThreadPool{
		finished:     false,
		workerN:      workerN,
		messageQueue: make(chan TPMessage, queueSize),
	}
	tp.initWorkers()
	return tp
}

// Start worker gorountines for this thread pool
func (tp *ThreadPool) initWorkers() {
	for i := 0; i < tp.workerN; i++ {
		worker := newWorker(tp.messageQueue)
		worker.init()
	}
}

// Put a new TPJob in message queue to be executed.
//
// It PANICS: If called after thread pool has finished
func (tp *ThreadPool) Run(j TPJob) {
	if tp.finished {
		panic("Use of thread pool after finishing it")
	}
	tp.messageQueue <- j
}

// Wait for all TPJobs to finish and terminate goroutines
func (tp *ThreadPool) Finish() {
	// Should receive a message for each goroutine that it has finished
	exitSignal := make(chan bool, tp.workerN)
	// Send "I'm finished" to main thread
	var exitMessage exitMessage = func() {
		exitSignal <- true
	}
	// Message all workers to stop
	for i := 0; i < tp.workerN; i++ {
		tp.messageQueue <- exitMessage
	}
	// Await acknowledge
	for i := 0; i < tp.workerN; i++ {
		<-exitSignal
	}
	tp.finished = true
}

// Single goroutine wrapper for Thread Pool
type worker struct {
	messageQueue chan TPMessage
}

func newWorker(messageQueue chan TPMessage) *worker {
	return &worker{
		messageQueue: messageQueue,
	}
}

// Initialize goroutine and wait for messages
func (w *worker) init() {
	go func() {
		for {
			message := <-w.messageQueue
			if shouldExit := message.execute(); shouldExit {
				return
			}
		}
	}()
}
