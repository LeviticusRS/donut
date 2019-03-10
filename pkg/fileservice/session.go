package fileservice

import (
    "errors"
    "github.com/sprinkle-it/donut/pkg/client"
)

type SessionFactory func(*client.Client, WorkerQueue) *Session

type SessionConfig struct {
    PriorityRequestCapacity int
    PassiveRequestCapacity  int
}

func (cfg SessionConfig) Build(cli *client.Client, workers WorkerQueue) *Session {
    return &Session{
        Client:   cli,
        priority: make(chan Request, cfg.PassiveRequestCapacity),
        passive:  make(chan Request, cfg.PassiveRequestCapacity),
        workers:  workers,
        done:     make(chan error, 1),
    }
}

type Session struct {
    *client.Client

    // The queue of priority requests that the session has received. Priority requests are for archives
    // that are critically needed in order for the client to function and should be handled before
    // passive requests.
    priority chan Request

    // The queue of passive requests that the session has received. Passive requests are for archives
    // that are not critically needed at the present time but will be needed eventually.
    passive chan Request

    workers WorkerQueue

    done chan error
}

func (s *Session) Process() {
    go func() {
        for {
            // First attempt to select from the priority requests so we can assure that the client is receiving files
            // that it immediately needs. If the session has been closed just return.
            select {
            case <-s.Quit():
                return
            case request := <-s.priority:
                s.submitRequest(request)
                continue
            default:
            }

            // If there was no priority request take either a request from the passive or priority queue.
            select {
            case <-s.Quit():
                return
            case request := <-s.priority:
                s.submitRequest(request)
            case request := <-s.passive:
                s.submitRequest(request)
            }
        }
    }()
}

func (s *Session) submitRequest(request Request) {
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

func (s *Session) enqueue(request Request, queue chan Request) {
    select {
    case queue <- request:
    default:
        s.Fatal(errors.New("fileservice: request queue is full"))
    }
}