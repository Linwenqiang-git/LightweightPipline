package models

type Response struct {
	Code    int
	Message string
	Data    interface{}
}

// 返回成功
func Success(data interface{}) Response {
	result := Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
	return result
}
func SuccessWithMessage(message string, data interface{}) Response {
	result := Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
	return result
}

// 返回失败
func Fail(code int, message string) Response {
	result := Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	return result
}
