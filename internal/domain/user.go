package domain

import "github.com/YourPainkiller/BHS_test/internal/dto"

type User struct {
	userId       int
	userName     string
	userPassword string
}

const MAXUSERNAMELEN = 12

func NewUser(userName, userPassword string) (*User, error) {
	user := User{}

	err := user.SetUsername(userName)
	if err != nil {
		return nil, err
	}

	err = user.SetPassword(userPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (u *User) SetUserId(userId int) error {
	if userId < 1 {
		return ErrInvalidUserId
	}
	u.userId = userId
	return nil
}

func (u *User) SetUsername(userName string) error {
	if len(userName) < 1 || len(userName) > MAXUSERNAMELEN {
		return ErrInvalidUsername
	}
	u.userName = userName
	return nil
}

func (u *User) SetPassword(userPassword string) error {
	u.userPassword = userPassword
	return nil
}

func (u *User) ToDTO() dto.UserDto {
	return dto.UserDto{
		UserId:       u.userId,
		UserName:     u.userName,
		UserPassword: u.userPassword,
	}
}
