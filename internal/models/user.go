package models

import "time"

//user представляет пользователя в системе
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

//CreateUserRequest - данные для создания пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//данные для входа
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//ответ с JWT токеном
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
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
	if v.Password == "" {
		return ErrPasswordRequired
	}
	if len(v.Password) < 6 {
		return ErrPasswordTooShort
	}
	return nil
}

//Validate проверяет LoginRequest
func (r *LoginRequest) Validate() error {
	if r.Username == "" {
		return ErrUsernameRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}
