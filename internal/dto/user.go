package dto

type UserDto struct {
	UserId       int    `json:"userId" db:"id"`
	UserName     string `json:"userName" db:"username"`
	UserPassword string `json:"userPassword" db:"password"`
}

type RegisterUserDto struct {
	UserName     string `json:"userName"`
	UserPassword string `json:"userPassword"`
}
