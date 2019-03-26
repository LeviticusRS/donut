package message

import "github.com/sprinkle-it/donut/pkg/buffer"

type StreamEncoder struct {
    buffer []byte
}

func NewStreamEncoder(capacity int) StreamEncoder {
    return StreamEncoder{buffer: make([]byte, capacity)}
}

func (e *StreamEncoder) Encode(msg Outbound, output *buffer.RingBuffer) error {
    buf := buffer.ByteBuffer{Bytes:e.buffer}

    if err := buf.PutUint8(uint8(msg.Config().Id)); err != nil {
        return err
    }

    // Mark where we start beginning writing the message so that we can determine the length of the message.
    start := buf.Offset

    // Skip over where we need to write the length of the packet. The number of bytes that dictate the length
    // is determined by the size of the packet.
    size := msg.Config().Size
    if err := buf.Skip(size.encodedLength()); err != nil {
        return err
    }

    if err := msg.Encode(&buf); err != nil {
        return err
    }

    end := buf.Offset

    length := end - start - size.encodedLength()

    // Move back to the start of the message and write the length.
    buf.Offset = start

    switch size {
    case SizeVariableByte:
        if err := buf.PutUint8(uint8(length)); err != nil {
            return err
        }
    case SizeVariableShort:
        if err := buf.PutUint16(uint16(length)); err != nil {
            return err
        }
    default:
        // Statically sized messages do not need their length to be encoded.
    }

    // Write the message bytes to the buffer.
    if _, err := output.Write(e.buffer[:end]); err != nil {
        return err
    }

    return nil
}
