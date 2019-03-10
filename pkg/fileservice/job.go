package fileservice

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/client"
)

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

    buf := buffer.ByteBuffer{Bytes:worker.buffer[:]}

    _ = buf.PutUint8(j.request.Index)
    _ = buf.PutUint16(j.request.Id)

    if err := writer.Write(worker.buffer[:3]); err != nil {
        j.done <- err
        return
    }

    offset := 0

    for offset < len(bytes) {
        if offset > 0 {
            if err := writer.Write(endOfChunk); err != nil {
                j.done <- err
                return
            }
        }

        write := len(bytes) - offset
        if write > chunkLength {
            write = chunkLength
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
