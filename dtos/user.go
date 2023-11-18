package dtos

type UserRegister struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"johndoe@mail.com"`
	Password string `json:"password" binding:"required,min=6" example:"JohnDoe123"`
}

type RegisterResponse struct {
	ID string `json:"user_id"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email" example:"johndoe@mail.com"`
	Password string `json:"password" binding:"required,min=6" example:"JohnDoe123"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserUpdateRequest struct {
	Username    string `json:"username" example:"johndoe123"`
	Email       string `json:"email" binding:"omitempty,email" example:"johndoe@mail.org"`
	Password    string `json:"password" binding:"required,min=6" example:"JohnDoe123"`
	NewPassword string `json:"new_password" binding:"omitempty,min=6"`
}

type UserResponse struct {
	Usename string `json:"username"`
}
