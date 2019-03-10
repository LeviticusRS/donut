package fileservice

import (
    "errors"
    "github.com/sprinkle-it/donut/pkg/client"
)

type SessionFactory func(*client.Client, WorkerPool) *Session

type SessionConfig struct {
    PriorityRequestCapacity int
    PassiveRequestCapacity  int
}

func (cfg SessionConfig) Build(cli *client.Client, workers WorkerPool) *Session {
    return &Session{
        Client:   cli,
        workers:  workers,
        priority: make(chan Request, cfg.PassiveRequestCapacity),
        passive:  make(chan Request, cfg.PassiveRequestCapacity),
        done:     make(chan error, 1),
    }
}

// A file service session serves files back to a client as it requests them. There are two types of requests a client
// can make, priority requests and passive requests. Priority requests are expected to be served before passive
// requests.
type Session struct {
    *client.Client

    // The pool of workers that the session can utilize to submit jobs.
    workers WorkerPool

    // The queue of priority requests that the session has received. Priority requests are for archives that are
    // critically needed in order for the client to function and should be handled before passive requests.
    priority chan Request

    // The queue of passive requests that the session has received. Passive requests are for archives that are not
    // critically needed at the present time but will be needed eventually.
    passive chan Request

    // Channel used to signal when a job has been completed.
    done chan error
}

func (s *Session) Process() {
    go func() {
        for {
            // Attempt to select from the priority requests so we can assure that the client is receiving files that it
            // immediately needs. If the session has been closed just return.
            select {
            case <-s.Quit():
                return
            case request := <-s.priority:
                s.submit(request)
                continue
            default:
                // No priority requests currently.
            }

            // If there was no priority request take either a request from the passive or priority queue.
            select {
            case <-s.Quit():
                return
            case request := <-s.priority:
                s.submit(request)
            case request := <-s.passive:
                s.submit(request)
            }
        }
    }()
}

// Enqueues a request to the priority queue.
func (s *Session) EnqueuePriority(request Request) {
    s.enqueue(request, s.priority)
}

// Enqueues a request to the passive queue.
func (s *Session) EnqueuePassive(request Request) {
    s.enqueue(request, s.passive)
}

// Enqueues a request to the provided channel. If the channel cannot immediately accept the request then the session
// will be closed with a fatal error.
func (s *Session) enqueue(request Request, queue chan Request) {
    select {
    case queue <- request:
    default:
        s.Fatal(errors.New("fileservice: request queue is full"))
    }
}

// Submits a request to a worker.
func (s *Session) submit(request Request) {
    select {
    case <-s.Quit():
        return
    case worker := <-s.workers:
        worker <- NewJob(s, request)

        select {
        case err := <-s.done:
            if err != nil && err != client.ErrClosed {
                s.Fatal(err)
            }
        }
    }
}
