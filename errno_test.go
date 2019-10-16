package xerrors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrno(t *testing.T) {
	e := Errno(1)
	var e2 error = Errno(1)
	var e3 error = Errno(2)

	//fmt.Printf("is = %v\n", errors.Is(e, e3))
	assert.False(t, errors.Is(e, e3))
	assert.True(t, errors.Is(e, e2))
}
