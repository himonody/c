package dto

type RegisterRequest struct {
	Username   string `json:"username" validate:"required,min=6,max=12"`   // 6-12 characters
	Password   string `json:"password" validate:"required,min=6,max=12"`   // 6-12 characters
	FriendCode string `json:"friendCode"`                                  // Optional invite code
	// Nickname is auto-generated, not provided by user
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

type UserInfo struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	RefCode      string `json:"refCode"`
	Status       string `json:"status"`
	OnlineStatus string `json:"onlineStatus"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,min=6,max=12"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=12"`
}

// SetPayPasswordRequest 设置/修改支付密码请求
type SetPayPasswordRequest struct {
	PayPassword    string `json:"payPassword" validate:"required,min=6,max=6"`    // 6位数字
	OldPayPassword string `json:"oldPayPassword,omitempty"`                       // 修改时需要旧密码
}

// UpdateProfileRequest 修改用户资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname,omitempty"` // 昵称
	Avatar   string `json:"avatar,omitempty"`   // 头像URL
}
