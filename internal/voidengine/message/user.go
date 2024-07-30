package message

type ListUserRequest struct {
	ListQueryBase
}

type ListUserResponse struct {
	ListResponseBase
}

type AddUserRequest struct {
	Username string `json:"username"    binding:"required,max=20"`       // 用户名
	Password string `json:"password"    binding:"required,min=8,max=20"` // 密码
}
