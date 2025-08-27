package herrors

import (
	"errors"
	"reflect"
)

const (
	DefaultServerErrorCode = 500
	DefaultParameterError  = 400
	DefaultAuthError       = 401
)

// NewServerError 运行时错误， 比如查询出错等， 定义时传入错误理由，使用时传入具体错误信息
func NewServerError(reason string) func(error) *HError {
	return func(err error) *HError {
		return &HError{Code: DefaultServerErrorCode, Reason: reason, DefMessage: err.Error(), BusinessError: err}
	}
}

// NewServerDefMessageError 运行时错误， 比如查询出错等， 定义时传入错误理由，使用时传入具体错误信息
func NewServerDefMessageError(reason, msg string) func(error) *HError {
	return func(err error) *HError {
		return &HError{Code: DefaultServerErrorCode, Reason: reason, DefMessage: msg, BusinessError: err}
	}
}

// NewBusinessServerError 自定义业务错误，比如用户重复等，定义错误时直接传入错误理由即可
func NewBusinessServerError(reason string) *HError {
	return &HError{Code: DefaultServerErrorCode, Reason: reason, DefMessage: reason, BusinessError: errors.New(reason)}
}

func NewParameterError(reason string) func(error) Herr {
	return func(err error) *HError {
		if err == nil {
			return &HError{Code: DefaultParameterError, Reason: reason, DefMessage: reason, BusinessError: errors.New(reason)}
		}
		return &HError{Code: DefaultParameterError, Reason: reason, DefMessage: err.Error(), BusinessError: err}
	}
}
func NewParameterHError(err error) *HError {
	return &HError{Code: DefaultParameterError, Reason: ReqParameterError, DefMessage: ReqParameterError, BusinessError: err}
}
func NewServerHError(err error) Herr {
	return &HError{Code: DefaultServerErrorCode, Reason: "ServerError", DefMessage: err.Error(), BusinessError: err}
}

func NewBadReqError(reason string) Herr {
	return &HError{Code: DefaultParameterError, Reason: reason, DefMessage: reason, BusinessError: errors.New(reason)}
}
func NewBadReqHError(err error) Herr {
	return &HError{Code: DefaultParameterError, Reason: ReqParameterError, DefMessage: ReqParameterError, BusinessError: err}
}
func NewAsServerError(err error, def *HError) *HError {
	if herr, is := IsHServerError(err); is {
		return herr
	}

	if def != nil {
		return def
	}

	return &HError{Code: DefaultServerErrorCode, Reason: ReasonStatusInternalHError, DefMessage: "Server Internal Error", BusinessError: err}
}

func IsHServerError(err error) (*HError, bool) {
	errType := reflect.TypeOf(err).String()
	if errType == "*HError" {
		return err.(*HError), true
	}
	return nil, false
}

var (
	BaseServerError           = New(DefaultServerErrorCode, "ServerError", "SERVER_ERROR")      // 服务器错误
	BaseParameterError        = New(DefaultParameterError, "ParameterError", "PARAMETER_ERROR") // 参数错误
	BaseFrequentRequestsError = NewBusinessServerError("BaseFrequentRequestsError")             // 请求过于频繁,请稍后再试

	BaseTokenEmpty      = New(DefaultAuthError, "TokenEmpty", "TOKEN_EMPTY_ERROR")            // token 不存在
	BaseTokenVerifyFail = New(DefaultAuthError, "TokenVerifyFail", "TOKEN_VERIFY_FAIL_ERROR") // token 失效
	SystemMaintenance   = NewBusinessServerError("SystemMaintenance")                         //系统维护中
	QueryFail           = NewServerError("QueryFail")                                         // 查询失败
	CreateFail          = NewServerError("CreateFail")                                        // 创建失败
	UpdateFail          = NewServerError("UpdateFail")                                        // 更新失败
	DeleteFail          = NewServerError("DeleteFail")                                        // 删除失败
	OperateFail         = NewServerError("OperateFail")                                       // 操作失败
	NotFound            = NewServerError("NotFound")                                          // 未找到
	DataIsExist         = NewBusinessServerError("DataIsExist")                               //数据已经存在
	ErrRecordNotFount   = NewBusinessServerError("ErrRecordNotFount")                         //记录未找到
	PleaseDoNotResubmit = NewBusinessServerError("PleaseDoNotResubmit")                       //请勿重复提交

	// 系统配置相关错误
	SysConfigurationNotExist = NewBusinessServerError("SysConfigurationNotExist") // 系统配置不存在

	// 导出相关错误
	ExportExcelError       = NewServerError("ExportExcelError")               // 导出Excel失败
	ExportTypeNotSupported = NewBusinessServerError("ExportTypeNotSupported") // 不支持的导出类型
	ExportColumnsEmpty     = NewBusinessServerError("ExportColumnsEmpty")     // 导出列配置为空
	ExportDataEmpty        = NewBusinessServerError("ExportDataEmpty")        // 导出数据为空
	ExportFileNameInvalid  = NewBusinessServerError("ExportFileNameInvalid")  // 导出文件名无效
	ExportWriteHeaderError = NewServerError("ExportWriteHeaderError")         // 写入表头失败
	ExportWriteDataError   = NewServerError("ExportWriteDataError")           // 写入数据失败
	ExportSaveFileError    = NewServerError("ExportSaveFileError")            // 保存文件失败
)
