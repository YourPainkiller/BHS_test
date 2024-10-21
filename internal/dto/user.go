package dto

type UserDto struct {
	UserId       int    `json:"id" db:"id"`
	UserName     string `json:"username" db:"username"`
	UserPassword string `json:"password" db:"password"`
}

type AuthUserDto struct {
	UserName     string `json:"username"`
	UserPassword string `json:"password"`
}
