package contract

type ListUserRequest struct {
	ListQuery
	Username *string `form:"username" json:"username" binding:"omitempty,max=32"`
}

type ListUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type AddUserRequest struct {
	Username   string `json:"username" binding:"required,max=32"`
	Password   string `json:"password" binding:"required,min=8,max=20"`
	RePassword string `json:"rePassword" binding:"required,min=8,max=20,eqfield=Password"`
	Email      string `json:"email" binding:"min=5,max=100"`
	Phone      string `json:"phone" binding:"min=4,max=20"`
}

type ModifyUserRequest struct {
	Username   string  `json:"username" binding:"required,max=32"`
	Password   *string `json:"password" binding:"min=8,max=20"`
	RePassword *string `json:"rePassword" binding:"min=8,max=20,eqfield=Password"`
	Email      *string `json:"email" binding:"min=5,max=100"`
	Phone      *string `json:"phone" binding:"min=4,max=20"`
}

type DeleteUserRequest struct {
	Username string `json:"username" binding:"required,max=32"`
}
