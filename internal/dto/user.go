package dto

type RegisterUserDto struct {
	UserPassword string `json:"password" db:"password" example:"yourpassword"`
	UserName     string `json:"username" db:"username" example:"yourusername"`
}

type UserDto struct {
	UserId       int    `json:"id" db:"id"`
	UserName     string `json:"username" db:"username"`
	UserPassword string `json:"password" db:"password"`
}

// model for swagger
type UserCredentials struct {
	UserPassword string `json:"password" example:"yourpassword"`
	UserName     string `json:"username" example:"yourusername"`
}
