package model

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// 错误码定义
const (
	CodeSuccess          = 0
	CodeInvalidParam     = 40001
	CodeUnauthorized     = 40101
	CodeTokenInvalid   = 40102
	CodeTokenExpired   = 40103
	CodeForbidden        = 40301
	CodeNotFound         = 40401
	CodeInternalError    = 50001
)
