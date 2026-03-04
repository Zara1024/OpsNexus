package errors

// AppError 定义统一应用错误结构。
type AppError struct {
	Code     int    `json:"code"`
	HTTPCode int    `json:"-"`
	Message  string `json:"message"`
	Detail   string `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrUnauthorized = &AppError{Code: 10001, HTTPCode: 401, Message: "未登录或登录已过期"}
	ErrForbidden    = &AppError{Code: 10003, HTTPCode: 403, Message: "无权限访问"}
	ErrNotFound     = &AppError{Code: 10004, HTTPCode: 404, Message: "资源不存在"}
	ErrValidation   = &AppError{Code: 10005, HTTPCode: 400, Message: "参数校验失败"}
	ErrInternal     = &AppError{Code: 10006, HTTPCode: 500, Message: "服务器内部错误"}
)
