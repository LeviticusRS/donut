package message

import (
    "fmt"
    "github.com/sprinkle-it/donut/pkg/buffer"
)

type StreamDecoderState int

const (
    DecodeIdentifier = 0
    DecodeLength     = 1
    AwaitBytes       = 2
)

type StreamDecoder struct {
    messages        DescriptorSet
    receivedMessage Descriptor
    receivedLength  int
    state           StreamDecoderState
    buffer          []byte
}

func NewStreamDecoder(messages DescriptorSet, capacity int) StreamDecoder {
    return StreamDecoder{
        messages: messages,
        state:    DecodeIdentifier,
        buffer:   make([]byte, capacity),
    }
}

func (d *StreamDecoder) Decode(r buffer.Readable) (Message, error) {
    switch d.state {
    case DecodeIdentifier:
        if !buffer.HasReadable(r) {
            return nil, nil
        }

        if _, err := r.Read(d.buffer[:1]); err != nil {
            return nil, err
        }

        id := d.buffer[0]

        config, ok := d.messages[id]
        if !ok {
            return nil, fmt.Errorf("message: unrecognized message %d", id)
        }

        d.state = DecodeLength
        d.receivedMessage = config
        fallthrough
    case DecodeLength:
        switch d.receivedMessage.Size {
        case SizeVariableByte:
            if !buffer.IsReadable(r, 1) {
                return nil, nil
            }

            if _, err := r.Read(d.buffer[:1]); err != nil {
                return nil, err
            }

            d.receivedLength = int(d.buffer[0])
        case SizeVariableShort:
            if !buffer.IsReadable(r, 2) {
                return nil, nil
            }

            if _, err := r.Read(d.buffer[:2]); err != nil {
                return nil, err
            }

            d.receivedLength = int(uint16(d.buffer[0])<<8 | uint16(d.buffer[1]))
        default:
            d.receivedLength = int(d.receivedMessage.Size)
        }

        d.state = AwaitBytes
        fallthrough
    case AwaitBytes:
        if !buffer.IsReadable(r, d.receivedLength) {
            return nil, nil
        }

        if _, err := r.Read(d.buffer[:d.receivedLength]); err != nil {
            return nil, err
        }

        msg := New(d.receivedMessage)
        buf := buffer.ByteBuffer{Bytes:d.buffer[:d.receivedLength]}

        if err := msg.Decode(&buf, d.receivedLength); err != nil {
            return nil, err
        }

        d.state = DecodeIdentifier
        return msg, nil
    default:
        panic("Unexpected packet decoder state")
    }
}
