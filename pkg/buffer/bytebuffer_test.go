package buffer

import (
	"strings"
	"testing"
)

func TestByteBuffer_GetCString(t *testing.T) {
	const ExpectedStringValue = "praise sino"

	buffer := NewByteBuffer(32)
	for i := 0; i < len(ExpectedStringValue); i++ {
		_ = buffer.PutUint8(ExpectedStringValue[i])
	}

	buffer.PutUint8(NullTerminator)
	buffer.Offset = 0

	readResult, err := buffer.GetCString()
	if err != nil {
		t.Error(err)
	}

	if readResult != ExpectedStringValue {
		t.Error("value mismatch: expected %i to match %i", readResult, ExpectedStringValue)
	}
}

func TestByteBuffer_PutCString(t *testing.T) {
	const ExpectedStringValue = "sini is a noob"

	buffer := NewByteBuffer(32)
	_ = buffer.PutCString(ExpectedStringValue)
	buffer.Offset = 0

	var bldr strings.Builder
	for i := 0; i < len(ExpectedStringValue)+1; i++ {
		value, _ := buffer.GetUint8()
		if value == 0 {
			break
		}

		bldr.WriteByte(value)
	}

	readResult := bldr.String()
	if readResult != ExpectedStringValue {
		t.Error("value mismatch: expected %i to match %i", readResult, ExpectedStringValue)
	}
}
