package message

type ListUserRequest struct {
	ListQueryBase
	Username *string `form:"username"	json:"username"		binding:"max=32"` // 用户名
}

type ListUserResponse struct {
	ListResponseBase
}

type AddUserRequest struct {
	Username   string `json:"username"	binding:"required,max=32"`       // 用户名
	Password   string `json:"password"	binding:"required,min=8,max=20"` // 密码
	RePassword string `json:"rePassword"	binding:"required,min=8,max=20,eqfield=Password"`
	Email      string `json:"email"		binding:"min=5,max=100"`
	Phone      string `json:"phone"		binding:"min=4,max=20"`
}

type ModifyUserRequest struct {
	Username   string  `json:"username"	binding:"required,max=32"` // 用户名
	Password   *string `json:"password"	binding:"min=8,max=20"`    // 密码
	RePassword *string `json:"rePassword"	binding:"min=8,max=20,eqfield=Password"`
	Email      *string `json:"email"	binding:"min=5,max=100"`
	Phone      *string `json:"phone"	binding:"min=4,max=20"`
}

type DeleteUserRequest struct {
	Username string `json:"username"	binding:"required,max=32"` // 用户名
}
