package gateway

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBuffer(t *testing.T) {
	w := make([]byte, 54)
	b := bytes.NewBuffer([]byte{})
	b.Grow(4096)
	b.Write(w)
	fmt.Printf("a--:%d,%d\n", b.Cap(), b.Len())
	b.Reset()
	fmt.Printf("b--:%d, %d\n", b.Cap(), b.Len())
}
