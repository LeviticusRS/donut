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

type Config struct {
    Id   uint8
    Size Size
    New  func() Message
}
