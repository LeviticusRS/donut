package file

type WorkerPool chan JobQueue

type Worker struct {
    pool            WorkerPool
    jobs            JobQueue
    archiveProvider ArchiveProvider
    buffer          [8]byte
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
