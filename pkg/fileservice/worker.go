package fileservice

var (
    endOfChunk  = []byte{255}
    chunkLength = 2048
)

type WorkerPool chan JobQueue

type Worker struct {
    pool     WorkerPool
    jobs     JobQueue
    provider ArchiveProvider
    buffer   [8]byte
}

func (w *Worker) Process() {
    go func() {
        for {
            w.pool <- w.jobs
            select {
            case job := <-w.jobs:
                job.Execute(w)
            }
        }
    }()
}
