package uerrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

// 本包基于go1.13设计 (go version >= go1.13)
// 声明:error工具包，不涉及任何业务函数
// API还处于摸索阶段，随时会变
var emptyError = errors.New("nil")

// ErrorMsg定义
type ErrorMsg struct {
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 返回给客户端的错误
	err     error  //保存内部错误，这个错误可能包装过很过次
}

// 把code, msg 实例化为新ErrorMsg类型
func NewMsg(code int, msg string) *ErrorMsg {
	return &ErrorMsg{ErrCode: code, ErrMsg: msg, err: errors.New(msg)}
}

// 把code, err 实例化成新的ErrorMsg类型
func New(code int, err error) *ErrorMsg {
	return &ErrorMsg{ErrCode: code, ErrMsg: err.Error(), err: err}
}

// 返回code
func (x *ErrorMsg) Code() int {
	return x.ErrCode
}

// 返回Error
func (x *ErrorMsg) Error() string {
	return x.err.Error()
}

// 给errors.Unwrap函数调用,用于解包错误
func (x *ErrorMsg) Unwrap() error {
	return x.err
}

// 把ErrorMsg 换成json 字符串
func (x *ErrorMsg) JsonString() string {
	all, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}

	return *(*string)(unsafe.Pointer(&all))
}

// 提取错误吗
func ToCode(err error) int {
	return ToErrorMsg(err).Code()
}

// 从错误链中找出ErrorMsg类型
// 找不到构造一个默认值
func ToErrorMsg(err error) *ErrorMsg {
	var e *ErrorMsg

	if err == nil {
		err = emptyError
		goto next
	}

	if errors.As(err, &e) {
		return e
	}

next:
	return &ErrorMsg{ErrCode: 0xff, ErrMsg: fmt.Sprintf("%s:No ErrorMsg type found", err)}
}

// 从错误里面检查ErrorMsg结构体，如果已经有，就不装包
// 没有就包装为XErrMsg结构
func Errorf(code int, format string, args ...interface{}) error {
	var x *ErrorMsg
	err := fmt.Errorf(fmt.Sprintf("code:%s", format), args...)

	if errors.As(err, &x) {
		return err
	}

	return New(code, err)
}
