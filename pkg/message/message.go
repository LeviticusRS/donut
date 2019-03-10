package message

import (
    "errors"
    "github.com/sprinkle-it/donut/pkg/buffer"
)

var (
    ErrDecodeNotSupported = errors.New("message: decoding is not supported for this message")
    ErrEncodeNotSupported = errors.New("message: encoding is not supported for this message")
)

// Type declaration for message descriptor sets which map a message descriptor by its id.
type DescriptorSet map[uint8]Descriptor

// Type declaration for functions that that provide messages of a certain type. The provider can give messages that are
// already created to be reused or create new messages.
type Provider func() Message

// Creates a new provider that provides a singleton value.
func ProvideSingleton(msg Message) Provider { return func() Message { return msg } }

type Message interface {
    // Gets the descriptor which specifies the id, size, and factory for creating new messages of this type.
    Descriptor() Descriptor

    // Decodes the message from the buffer. It is expected that the buffer contains the number of bytes required to
    // decode the message, if not then this function will return an error. If this function returns no error then
    // decoding the message was successful. Implementations of this method are expected not to validate the decoded data
    // to any contextual source. Exceptions to this include where there are checks to assure that the message is valid.
    // Any other validation should be done outside of the message by receivers of the message.
    Decode(buf *buffer.ByteBuffer, length int) error

    // Encodes the message to a buffer. It is expected that the buffer capacity is enough to contain the entire encoded
    // message, if not then this function will return an error.
    Encode(*buffer.ByteBuffer) error
}

// Creates a new message from the provided descriptor.
func New(descriptor Descriptor) Message {
    return descriptor.Provider()
}
