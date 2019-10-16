package xerrors

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

// XErrorMsg定义
type XErrorMsg struct {
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 返回给客户端的错误
	err     error  //保存内部错误，这个错误可能包装过很过次
}

// 把code, msg 实例化为新XErrorMsg类型
func NewMsg(code int, msg string) *XErrorMsg {
	return &XErrorMsg{ErrCode: code, ErrMsg: msg, err: errors.New(msg)}
}

// 把code, err 实例化成新的XErrorMsg类型
func New(code int, err error) *XErrorMsg {
	return &XErrorMsg{ErrCode: code, ErrMsg: err.Error(), err: err}
}

// 返回code
func (x *XErrorMsg) Code() int {
	return x.ErrCode
}

// 返回Error
func (x *XErrorMsg) Error() string {
	return x.err.Error()
}

// 给errors.Unwrap函数调用,用于解包错误
func (x *XErrorMsg) Unwrap() error {
	return x.err
}

// 把XErrorMsg 换成json 字符串
func (x *XErrorMsg) JsonString() string {
	all, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}

	return *(*string)(unsafe.Pointer(&all))
}

// 提取错误吗
func ToCode(err error) int {
	return ToXErrorMsg(err).Code()
}

// 从错误链中找出XErrorMsg类型
// 找不到构造一个默认值
func ToXErrorMsg(err error) *XErrorMsg {
	var e *XErrorMsg

	if err == nil {
		err = emptyError
		goto next
	}

	if errors.As(err, &e) {
		return e
	}

next:
	return &XErrorMsg{ErrCode: 0xff, ErrMsg: fmt.Sprintf("%s:No XErrorMsg type found", err)}
}

// 从错误里面检查XErrorMsg结构体，如果已经有，就不装包
// 没有就包装为XErrMsg结构
func Errorf(code int, format string, args ...interface{}) error {
	var x *XErrorMsg
	err := fmt.Errorf(fmt.Sprintf("code:%s", format), args...)

	if errors.As(err, &x) {
		return err
	}

	return New(code, err)
}
