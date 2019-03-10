package fileservice

var (
    endOfChunk  = []byte{255}
    chunkLength = 2048
)

type WorkerQueue chan JobQueue

type JobQueue chan Job

type Worker struct {
    pool     WorkerQueue
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
