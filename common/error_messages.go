package common

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
	ErrInvalidToken = errors.New("invalid token")
)

var listErrorData = []errData{
	// Define your error data here, for example:
	{Code: "user_not_found", HTTPCode: 404, MessageViVn: "Không tìm thấy người dùng", MessageEnUs: "User not found"},
	{Code: "invalid_token", HTTPCode: 401, MessageViVn: "Token không hợp lệ", MessageEnUs: "Invalid token"},
	// Add more error codes and their messages here...
}

var (
	AllErrors *MasterErrData
)

// FetchMasterErrData initializes AllErrors with the error data.
func FetchMasterErrData() {
	AllErrors = NewMasterErrData()
	AllErrors.fetchAll()
}

type errData struct {
	Code        string `json:"code" gorm:"column:code"`
	HTTPCode    int    `json:"http_code" gorm:"column:http_code"`
	MessageViVn string `json:"message_vi_vn" gorm:"column:message_vi_vn"`
	MessageEnUs string `json:"message_en_us" gorm:"column:message_en_us"`
}

type ExtraData struct {
	OrderID int64 `json:"order_id,omitempty"`
}

type LocalizeErrRes struct {
	Code      string     `json:"code,omitempty"`
	Message   string     `json:"message,omitempty"`
	HTTPCode  int        `json:"-"`
	Internal  string     `json:"internal,omitempty"`
	ExtraData *ExtraData `json:"extra_data,omitempty"`
}

func (a *LocalizeErrRes) Error() string {
	return a.Code
}

type MasterErrData struct {
	mutex sync.Mutex
	data  map[string]errData
}

// Initialize MasterErrData
func NewMasterErrData() *MasterErrData {
	return &MasterErrData{
		data: make(map[string]errData), // Initialize the map
	}
}

// Load all error data into the MasterErrData structure
func (a *MasterErrData) fetchAll() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, errMessage := range listErrorData {
		a.data[errMessage.Code] = errMessage
	}
}

// New creates a new localized error response
func (a *MasterErrData) New(err error, language string, internal ...string) *LocalizeErrRes {
	errRes := new(LocalizeErrRes)
	ok := errors.As(err, &errRes)
	if !ok {
		errRes = &LocalizeErrRes{
			Code:    "bad_request",
			Message: "Đã có lỗi xảy ra, vui lòng thử lại!",
		}
		if len(internal) > 0 {
			errRes.Internal = internal[0]
		}
		errFromDB, exists := a.data[err.Error()]
		if exists {
			errRes.Code = errFromDB.Code
			errRes.HTTPCode = errFromDB.HTTPCode
			switch language {
			case "vi":
				errRes.Message = errFromDB.MessageViVn
			default:
				errRes.Message = errFromDB.MessageEnUs
			}
		} else {
			errRes.HTTPCode = 400
		}
	}

	if len(internal) > 0 {
		errRes.Internal = internal[0]
	}
	return errRes
}

func (a *LocalizeErrRes) SetMessage(message string) *LocalizeErrRes {
	a.Message = message
	return a
}

func (a *LocalizeErrRes) ReplaceDescByVars(args ...interface{}) *LocalizeErrRes {
	for _, arg := range args {
		a.Message = fmt.Sprintf(a.Message, arg)
	}
	return a
}

func (a *LocalizeErrRes) SetOrderIDToExtraData(orderID int64) *LocalizeErrRes {
	if a.ExtraData == nil {
		a.ExtraData = new(ExtraData)
	}
	a.ExtraData.OrderID = orderID
	return a
}

func (a *LocalizeErrRes) ConvertToBaseError() Response {
	res := BaseResponse(REQUEST_FAILED, a.Message, a.Internal, a.ExtraData)
	res.SetErrorCode(a.Code)
	return res
}

func AbortWithError(c *gin.Context, err error) {
	if AllErrors == nil {
		FetchMasterErrData()
	}
	errJSON := AllErrors.New(err, "vi", err.Error())
	c.AbortWithStatusJSON(errJSON.HTTPCode, errJSON.ConvertToBaseError())
}
