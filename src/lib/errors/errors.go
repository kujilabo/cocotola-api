package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

var ErrorfFunc = fmt.Errorf

func init() {

}

func UseFmtErrorf() {
	ErrorfFunc = fmt.Errorf
}

func UseXerrorsErrorf() {
	ErrorfFunc = xerrors.Errorf
}

func Errorf(format string, a ...interface{}) error {
	return ErrorfFunc(format, a...)
}
