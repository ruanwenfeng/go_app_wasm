package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v6/pkg/log"
	"go.uber.org/zap/buffer"
)

type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

type Payload struct {
	// [redacted]
}

func (p Payload) UploadToS3() error {
	time.Sleep(1 * time.Second)
	fmt.Println("UploadToS3")
	return nil
}

var (
	MaxWorker = 10 //os.Getenv("MAX_WORKERS")
	MaxQueue  = 10 //os.Getenv("MAX_QUEUE")
)

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				if err := job.Payload.UploadToS3(); err != nil {
					log.Errorf("Error uploading to S3: %s", err.Error())
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {

	// if r.Method != "POST" {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }

	// Read the body into a string for json decoding
	var content = &PayloadCollection{}
	content.Payloads = append(content.Payloads, Payload{})
	// err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&content)
	// if err != nil {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// Go through each payload and queue items individually to be posted to S3
	for _, payload := range content.Payloads {

		// let's create a job with the payload
		work := Job{Payload: payload}

		// Push the work onto the queue.
		JobQueue <- work
	}

	w.WriteHeader(http.StatusOK)
}

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, maxWorkers: maxWorkers}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

func main() {

	// mmm := make(chan chan int, 4)
	// a := make(chan int)
	// b := make(chan int)
	// c := make(chan int)
	// d := make(chan int)
	// fmt.Println("======")

	// mmm <- a
	// fmt.Println("a")

	// mmm <- b
	// fmt.Println("b")

	// mmm <- c
	// fmt.Println("c")

	// mmm <- d
	// fmt.Println("d")

	// fmt.Println("-----")
	JobQueue = make(chan Job, 10)
	dispatcher := NewDispatcher(200)
	dispatcher.Run()
	http.HandleFunc("/postjob", payloadHandler)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		buf := &buffer.Buffer{}
		buf.AppendString("hello world")
		w.Write(buf.Bytes())
	})
	http.ListenAndServe(":8089", nil)
}
