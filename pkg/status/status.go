package status

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    okayDescriptor               = message.Descriptor{Id: 0, Size: 0, Provider: message.Provide(Okay)}
    unsupportedVersionDescriptor = message.Descriptor{Id: 6, Size: 0, Provider: message.Provide(UnsupportedVersion)}
    fullDescriptor               = message.Descriptor{Id: 7, Size: 0, Provider: message.Provide(Full)}

    Okay               = okay{}
    UnsupportedVersion = unsupportedVersion{}
    Full               = full{}
)

type okay struct{}

func (okay) Descriptor() message.Descriptor { return okayDescriptor }

func (okay) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (okay) Encode(b *buffer.ByteBuffer) error { return nil }

type unsupportedVersion struct{}

func (unsupportedVersion) Descriptor() message.Descriptor { return unsupportedVersionDescriptor }

func (unsupportedVersion) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (unsupportedVersion) Encode(b *buffer.ByteBuffer) error { return nil }

type full struct{}

func (full) Descriptor() message.Descriptor { return fullDescriptor }

func (full) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (full) Encode(b *buffer.ByteBuffer) error { return nil }
