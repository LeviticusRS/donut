package client

import (
    "errors"
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
    "go.uber.org/zap"
    "net"
    "sync"
    "time"
)

var (
    ErrNotActive = errors.New("client: this operation could not be performed because the client is not active")
    ErrClosed    = errors.New("client: this operation could not be performed because the client is closed")
)

// Creates a new client for the given connection and router.
type Factory func(net.Conn, *zap.Logger, MailRouter) *Client

// Type declaration for client callback functions.
type Callback func(*Client)

// Generates a identifier for a client. Implementations of this function type
// must generate an identifier that is at least unique for the application context
// that it is used for.
type IdentifierGenerator func() uint64

// Generates identifiers that start at a given value. This function should not be used by multiple threads.
func IncrementalGenerator(start uint64) IdentifierGenerator {
    counter := uint64(start)
    return func() uint64 { counter++; return counter }
}

// Flush write automatically flushes bytes that are written to the client once the byte counter reaches the output
// capacity. This implementation expects the clients output buffer is empty when first being used.
type FlushWriter struct {
    *Client
    counter int
}

func (w *FlushWriter) Write(b []byte) error {
    if w.counter+len(b) >= w.output.Capacity() {
        if err := w.Flush(); err != nil {
            return err
        }
        w.counter = 0
    }
    w.counter += len(b)
    return w.Client.Write(b)
}

// Type declaration for writer commands.
type outputCommand interface{}

// Writes bytes to the output buffer.
type writeBytes struct {
    bytes []byte
}

// Writes a message to the output buffer.
type writeMessage struct {
    msg message.Message
}

// Flushes the bytes from the output buffer to the connection.
type flushBytes struct{}

type Client struct {
    id uint64

    connection net.Conn

    logger *zap.Logger

    // All of the bytes read in from the connection will be written to this buffer. If the buffer cannot accept any more
    // bytes it will cause the client to close.
    input buffer.RingBuffer

    // All of the bytes being written to the client will be first written to this buffer. This allows callers to have
    // control over when bytes are flushed to the client.
    output buffer.RingBuffer

    // Blocking channel for output commands. When an operation wants to interact with the output buffer and encoder a
    // command needs to be published to this channel so that the operation can be executed synchronously.
    outputCommands chan outputCommand

    // Decodes byte streams into messages. Contains the state so that messages that have not been entirely transmitted
    // can be decoded when the bytes are received.
    decoder message.StreamDecoder

    // Encodes byte streams into messages. Encoding is stateless but in the future it might become stateful when a
    // cipher stream is needed to encode the message id for secure transport.
    encoder message.StreamEncoder

    // All of the received messages will be buffered to this channel. If this channel ever reaches its capacity the
    // client will close and an error will be reported.
    messages chan message.Message

    router MailRouter

    // Mutex which handles locking when operations need to check state that is not maintained by a go routine.
    mutex sync.Mutex

    // Flag for if the client processors have started. This will only be set once after the process function
    // has been called. Prevents multiple processors from fighting with resources.
    active bool

    // Flag for if the client is closed. This flag will be administered by the mutex to assure that it can
    // be accessed by multiple go routines querying about its state.
    closed bool

    // Signal for when the client was closed and therefore cannot handle any more operations. It is okay
    // for multiple threads to block on this channel because this channel will only ever be closed and
    // no values will be sent on it.
    quit chan struct{}
}

// Gets the client's identifier. This is assigned on creation and is guaranteed to be unique for the lifetime of the
// client. This identifier can be used to map other information pertaining to a client.
func (c *Client) Id() uint64 {
    return c.id
}

func (c *Client) RemoteAddress() net.Addr {
    return c.connection.RemoteAddr()
}

func (c *Client) Info(message string) {
    c.logger.Info(message, zap.Uint64("id", c.Id()), zap.Stringer("address", c.RemoteAddress()), )
}

// Writes bytes to the output buffer.
func (c *Client) Write(b []byte) error {
    if err := c.check(); err != nil {
        return err
    }

    c.outputCommands <- writeBytes{bytes: b}
    return nil
}

// Writes a message to the output buffer.
func (c *Client) Send(msg message.Message) error {
    if err := c.check(); err != nil {
        return err
    }

    c.outputCommands <- writeMessage{msg: msg}
    return nil
}

// Writes a message to the output buffer and flushes afterward.
func (c *Client) SendNow(msg message.Message) error {
    if err := c.check(); err != nil {
        return err
    }

    c.outputCommands <- writeMessage{msg: msg}
    c.outputCommands <- flushBytes{}
    return nil
}

// Flushes the output buffer to the connection.
func (c *Client) Flush() error {
    if err := c.check(); err != nil {
        return err
    }

    c.outputCommands <- flushBytes{}
    return nil
}

// Process received messages and publish them to the router.
func (c *Client) processMessages() {
    go func() {
        for {
            select {
            case <-c.quit:
                // Client was closed, stop dispatching messages.
                return
            case msg := <-c.messages:
                c.router.Publish(c, msg)
            }
        }
    }()
}

// Process reading from the connection and decoding messages.
func (c *Client) processInput() {
    go func() {
        transfer := make([]byte, c.input.Capacity())
        for {
            select {
            case <-c.quit:
                // Client was closed, stop reading from connection.
                return
            default:
                // Continue reading from connection.
            }

            count, err := c.connection.Read(transfer)
            if err != nil {
                if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
                    continue
                }
                c.Fatal(err)
                return
            }

            if _, err := c.input.Write(transfer[:count]); err != nil {
                c.Fatal(err)
                return
            }

            // Keep reading in bytes until the stream decoder says that it currently cannot decode a certain message.
            for buffer.HasReadable(&c.input) {
                msg, err := c.decoder.Decode(&c.input)
                if err != nil {
                    c.Fatal(err)
                    return
                }

                // Read bytes didn't contain a message and needs more bytes to continue.
                if msg == nil {
                    break
                }

                select {
                case c.messages <- msg:
                    // Successfully buffered message to client
                default:
                    // Failed to buffer message. Possibly because of misconfiguration, slow down, or denial of service.
                    c.Fatal(errors.New("client: failed to buffer message, message queue full"))
                    return
                }
            }
        }
    }()
}

// Process all of the output commands for the client.
func (c *Client) processOutput() {
    go func() {
        transfer := make([]byte, c.input.Capacity())
        for {
            select {
            case <-c.quit:
                // Client was closed, stop handling commands.
                return
            case cmd := <-c.outputCommands:
                switch cmd := cmd.(type) {
                case writeBytes:
                    if _, err := c.output.Write(cmd.bytes); err != nil {
                        c.Fatal(err)
                        return
                    }
                case writeMessage:
                    if err := c.encoder.Encode(cmd.msg, &c.output); err != nil {
                        c.Fatal(err)
                        return
                    }
                case flushBytes:
                    if !buffer.HasReadable(&c.output) {
                        return
                    }

                    count, err := c.output.Read(transfer[:c.output.Readable()])
                    if err != nil {
                        c.Fatal(err)
                        return
                    }

                    if _, err = c.connection.Write(transfer[:count]); err != nil {
                        c.Fatal(err)
                        return
                    }
                }
            }
        }
    }()
}

// Starts processing the input and output. This function is called when the server is ready to begin calling operations
// to the channel.
func (c *Client) Process() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if c.active {
        return
    }

    // The order in which these are called is important. Message processing is dependant on input processing. Output
    // processing has no dependants.
    c.processMessages()
    c.processInput()
    c.processOutput()

    c.active = true
}

func (c *Client) Fatal(err error) {
    if err := c.check(); err != nil {
        return
    }

    c.logger.Warn("Client encountered fatal error",
        zap.Uint64("id", c.id),
        zap.Stringer("address", c.connection.RemoteAddr()),
        zap.Error(err),
    )

    c.Close()
}

// Closes the client. This will alert the close subscribers by closing the quit channel. Returns if the client
// was closed, this function will return false if the client has already been closed.
func (c *Client) Close() bool {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if c.closed {
        return false
    }

    // Alert all of the subscribers listening for the client to close.
    close(c.quit)

    // Alert the reader and writer that the client is closing and that they should exit.
    _ = c.connection.SetDeadline(time.Now())

    // Close the connection gracefully.
    _ = c.connection.Close()

    return true
}

func (c *Client) Quit() <-chan struct{} {
    return c.quit
}

// Gets if the client is closed.
func (c *Client) Closed() bool {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    return c.closed
}

// Checks if the client is closed or started and return an error if either condition is true.
func (c *Client) check() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if c.closed {
        return ErrClosed
    }

    if c.active {
        return nil
    }

    return ErrNotActive
}

// Spawns a go routine that will wait for the client to close and then call the given function.
func (c *Client) OnClosed(callback Callback) {
    go func() {
        <-c.quit
        callback(c)
    }()
}
