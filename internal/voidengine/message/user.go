package message

type ListUserRequest struct {
	ListQueryBase
}

type ListUserResponse struct {
	ListResponseBase
}

type AddUserRequest struct {
	Username string `json:"username"    binding:"required,max=32"`       // 用户名
	Password string `json:"password"    binding:"required,min=8,max=20"` // 密码
}

type ModifyUserRequest struct {
	Username string `json:"username"    binding:"required,max=32"`       // 用户名
	Password string `json:"password"    binding:"required,min=8,max=20"` // 密码
}

type DeleteUserRequest struct {
	Username string `json:"username"    binding:"required,max=32"` // 用户名
}
