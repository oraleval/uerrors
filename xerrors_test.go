package xerrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	tErrRead  = 0x0001
	tErrWrite = 0x0002
)

func TestNew(t *testing.T) {
	// 测试New函数返回nil的情况
	errRead := New(tErrRead, errors.New("read fail"))
	assert.NotNil(t, errRead)
}

func TestNewMsg(t *testing.T) {
	// 测试NewMsg函数返回nil的情况
	errRead := NewMsg(tErrWrite, "write fail")
	assert.NotNil(t, errRead)
}

func TestCode(t *testing.T) {
	// 测试Code()函数，用于返回errcode
	errRead := NewMsg(tErrRead, "read fail")
	assert.Equal(t, errRead.Code(), tErrRead)
}

func TestError(t *testing.T) {
	// 测试Error()函数，用于返回error信息
	errError := NewMsg(tErrWrite, "write fail")
	assert.Equal(t, errError.Error(), "write fail")
}

func TestUnwrap(t *testing.T) {
	// 测试Unwrap函数，看是否拆包成功
	e1 := errors.New("e1")
	e2 := fmt.Errorf("%w", e1)

	e3 := New(tErrWrite, e2)

	e := errors.Unwrap(e3)
	e = errors.Unwrap(e)

	assert.Equal(t, e, e1)
}

func TestJsonString(t *testing.T) {
	// 测试JsonString 返回字符串，解码之后和原来的一样
	em1 := NewMsg(tErrRead, "read fail")
	str := em1.JsonString()

	em2 := XErrorMsg{}
	err := json.Unmarshal([]byte(str), &em2)
	assert.NoError(t, err)
	assert.Equal(t, em1.ErrCode, em2.ErrCode)
	assert.Equal(t, em1.ErrMsg, em2.ErrMsg)
}

func TestToCode(t *testing.T) {
	// 测试ToCode和传入是否不一样
	var e error = NewMsg(tErrRead, "test")
	e = fmt.Errorf("%w", e)

	assert.Equal(t, ToCode(e), tErrRead)
}

func TestToXerror(t *testing.T) {
	// 测试
	e1 := NewMsg(tErrWrite, "test")
	e := fmt.Errorf("aa :%w", e1)

	e2 := ToXErrorMsg(e)

	assert.Equal(t, e1, e2)
}

func TestErrorf(t *testing.T) {
	// 测试两种情况

	var err error = NewMsg(tErrRead, "new msg")
	// 情况1 error链有XErrorMsg
	err1 := Errorf(tErrWrite, "wrap %w", err)

	assert.Equal(t, ToXErrorMsg(err1), err.(*XErrorMsg))

	// 情况2 error链没有XErrorMsg
	err2 := Errorf(tErrWrite, "wrap %w", errors.New("test"))

	var x *XErrorMsg
	assert.True(t, errors.As(err2, &x))
}
