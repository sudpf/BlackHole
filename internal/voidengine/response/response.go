package response

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(code int, message string) *ApiResponse {
	return &ApiResponse{Code: code, Message: message}
}

func (r *ApiResponse) WithData(data interface{}) *ApiResponse {
	return &ApiResponse{Code: r.Code, Message: r.Message, Data: data}
}

/*
 * 0 - 10000 保留给通用的错误码
 * 每个模块同样保留100个错误码
 */
var (
	ApiSuccess   = &ApiResponse{Code: 0, Message: "Success"}
	InvalidParam = &ApiResponse{Code: 1, Message: "Invalid Param"}

	InvalidUserName = &ApiResponse{Code: 100001, Message: "Invalid UserName"}
	UserErrorEnd    = &ApiResponse{Code: 100100, Message: "User Error end"}
)
