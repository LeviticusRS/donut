package message

type Size int

const (
    SizeVariableByte  = -1
    SizeVariableShort = -2
)

func (s Size) EncodedLength() int {
    switch s {
    case SizeVariableByte:
        return 1
    case SizeVariableShort:
        return 2
    default:
        return 0
    }
}

func StaticSize(size int) Size {
    if size < 0 {
        panic("Size must be greater than or equal to zero")
    }
    return Size(size)
}

type Descriptor struct {
    Id       uint8
    Size     Size
    Provider Provider
}
