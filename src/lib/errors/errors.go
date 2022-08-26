package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

var ErrorfFunc = xerrors.Errorf

func init() {

}

func UseFmtErrorf() {
	ErrorfFunc = fmt.Errorf
}

func Errorf(format string, a ...interface{}) error {
	return ErrorfFunc(format, a...)
}
