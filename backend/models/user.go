package models

type User struct {
	ID         int    `json:"userId"`
	Username   string `json:"username"`
	Password   string `json:"-"` // 不返回给前端
	Followers  int    `json:"followers"`
	Following  int    `json:"following"`
	Avatar     string `json:"avatar"`
	DynamicNum int    `json:"dynamicNum"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"userResponse"`
}
type UserResponse struct {
	Username   string `json:"username"`
	Followers  int    `json:"followers"`
	Following  int    `json:"following"`
	Avatar     string `json:"avatar"`
	DynamicNum int    `json:"dynamicNum"`
	UserId     int    `json:"userId"`
	IsLogin    bool   `json:"isLogin"`
}
type SetupRequest struct {
	SetUsername string `json:"setusername" binding:"required"`
	SetPassword string `json:"setpassword" binding:"required"`
}
type SetupResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"userResponse"`
}

type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}
