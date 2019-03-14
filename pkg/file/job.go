package file

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/client"
)

var (
    // At the end of every chunk an acknowledge byte is written to the client so that it knows that the next chunk
    // is starting. If it does not receive this byte then the client drops the current request. This feature can
    // be used to halt a current request like in the case of the service receiving a priority request that should
    // be handled before the current request.
    EndOfChunk = []byte{255}

    // The length of a chunk in bytes.
    ChunkLength = 2048
)

// Type declaration for channels which are used to queue jobs for workers.
type JobQueue chan Job

// A job serves an archive to a session. Archives are served by first writing the archive we are writing, then
// chunks of the archive thereafter until the entire archive is written to the session.
type Job struct {
    session *Session
    request Request
    done    chan error
}

func NewJob(session *Session, request Request) Job {
    return Job{
        session: session,
        request: request,
        done:    session.done,
    }
}

func (j Job) Execute(worker *Worker) {
    bytes, err := worker.provider(j.request.Index, j.request.Id)
    if err != nil {
        j.done <- err
        return
    }

    writer := client.FlushWriter{Client: j.session.Client}

    buf := buffer.ByteBuffer{Bytes: worker.buffer[:]}

    _ = buf.PutUint8(j.request.Index)
    _ = buf.PutUint16(j.request.Id)

    if err := writer.Write(worker.buffer[:buf.Offset]); err != nil {
        j.done <- err
        return
    }

    offset := 0

    for offset < len(bytes) {
        if offset > 0 {
            if err := writer.Write(EndOfChunk); err != nil {
                j.done <- err
                return
            }
        }

        write := len(bytes) - offset
        if write > ChunkLength {
            write = ChunkLength
        }

        if err := writer.Write(bytes[offset : offset+write]); err != nil {
            j.done <- err
            return
        }

        offset += write
    }

    if err := j.session.Flush(); err != nil {
        j.done <- err
        return
    }

    j.done <- nil
}
