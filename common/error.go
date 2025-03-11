package common

import (
	"fmt"
	"strings"
)

const (
	TokenNotFound     = "TokenNotFound"
	TokenUnAuthorized = "TokenUnAuthorized"
	PerDenied         = "PerDenied"
	ResultFailed      = "ResultFailed"
	InternalError     = "InternalError"
)

const (
	DuplicatedDataErr = "Lỗi! Trùng lặp dữ liệu"
	DefaultError      = "Lỗi! Xảy ra lỗi không xác định, vui lòng liên hệ quản trị viên"
)

var DataIsNullErr = func(obj string) string {
	return fmt.Sprintf("%v cannot use nil", obj)
}

var DataIsExisted = func(obj string) string {
	return fmt.Sprintf("%v is existed", obj)
}

var DataIsSmallerZero = func(obj string) string {
	return fmt.Sprintf("%v is not smaller zero", obj)
}

var DataIsBeforeNow = func(obj string) string {
	return fmt.Sprintf("%v is not before now", obj)
}

var ErrorWrapper = func(prefix string, err error) error {
	return fmt.Errorf("%v: %v", prefix, err.Error())
}

var PgErrorTransform = func(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate key value") {
		return fmt.Errorf(DuplicatedDataErr)
	}

	return err
}
