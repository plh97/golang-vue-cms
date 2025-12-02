package v1

import "go-nunu/internal/model"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Image    string `json:"image"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" example:"alan"`
	Image    string `json:"image"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type GetUserListResponseData struct {
	// UserId   string `json:"userId"`
	// Nickname string `json:"nickname" example:"alan"`
	// Email    string `json:"email" example:"alan"`
	List []model.User `json:"list"`
}
type GetUserListResponse struct {
	Response
	Data GetUserListResponseData
}
