package dto

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Token     string        `json:"token"`
	AdminInfo AdminUserInfo `json:"adminInfo"`
}

type AdminUserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
	Role     int    `json:"role"`
	Status   int    `json:"status"`
}

type AdminLogoutResponse struct {
	Message string `json:"message"`
}
