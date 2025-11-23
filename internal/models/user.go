package models

import "time"

//user представляет пользователя в системе
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

//CreateUserRequest - данные для создания пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
}

//Validate - проверяет корректность данных
func (v *CreateUserRequest) Validate() error {
	if v.Username == "" {
		return ErrUsernameRequired
	}
	if len(v.Username) > 50 {
		return ErrUsernameTooLong
	}
	if len(v.Username) < 3 {
		return ErrUsernameTooShort
	}
	return nil
}
