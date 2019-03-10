package client

import (
    "errors"
    "github.com/sprinkle-it/donut/pkg/message"
)

// Type declaration for handlers that execute logic for received mail. Implementations of this type
// should never indefinitely block as it will cause go routines to leak.
type MailHandler func(Mail)

// A receiver is a wrapper to declare which messages a handler accepts.
type MailReceiver struct {
    Handler MailHandler
    Accept  []message.Descriptor
}

// Note(hadyn): The reason why messages and handlers are one to one are because if there are multiple
// receivers for a particular message then when a receiver for a message blocks the message may not be
// published to other receivers that accept the message. It is fairly expensive to spin up a go routine
// per message being sent as there could potentially be thousands or tens to hundreds of thousands per
// minute. This adds up quickly, the scheduler and garbage collector would be strained. The alternative
// to this implementation is to pass the message to a queue where a worker accepts a publish job. This has
// similar issues with locking out workers and in my opinion seems to be extraneous.
type MailRouter struct {
    handlers map[uint8]MailHandler
    accepted message.DescriptorSet
}

func NewMailRouter(receivers []MailReceiver) (MailRouter, error) {
    router := MailRouter{
        handlers: make(map[uint8]MailHandler),
        accepted: make(message.DescriptorSet),
    }

    for _, receiver := range receivers {
        for _, descriptor := range receiver.Accept {
            if _, ok := router.handlers[descriptor.Id]; ok {
                return MailRouter{}, errors.New("client: multiple receivers cannot accept the same message")
            }
            router.handlers[descriptor.Id] = receiver.Handler
            router.accepted[descriptor.Id] = descriptor
        }
    }

    return router, nil
}

func (r MailRouter) Publish(source *Client, msg message.Message) {
    if handler, ok := r.handlers[msg.Descriptor().Id]; ok {
        handler(Mail{Source: source, Message: msg})
    }
}

// A wrapper over a received message from a client.
type Mail struct {
    Source  *Client
    Message message.Message
}
