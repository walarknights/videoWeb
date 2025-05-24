package models

import "time"

type Video struct {
	ID        uint   `json:"vedioId"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Update    int64  `json:"update"`
	Likes     int    `json:"likes"`
	Favorites int    `json:"favorites"`
}
type Comment struct {
	ID         uint      `json:"commentId"`
	Content    string    `json:"content"`
	UserId     uint      `json:"userId"`
	UserName   string    `json:"userName"`
	UserAvatar string    `json:"userAvatar"`
	VideoID    uint      `json:"vedioId"`
	CreatedAt  time.Time `json:"createdAt"`
	ParentId   *uint     `json:"parentId"`
	RootId     uint      `json:"rootId"`
	ParentUser *UserCom  `gorm:"foreignKey:ParentUserID" json:"parentUser,omitempty"` // 被回复的用户信息
}

type CreateCommentInput struct {
	VideoID    uint   `json:"videoId" binding:"required"`
	UserID     uint   `json:"userId" binding:"required"`
	UserAvatar string `json:"userAvatar" binding:"required"`
	UserName   string `json:"userName" binding:"required"`
	Content    string `json:"content" binding:"required"`
	ParentID   *uint  `json:"parentId"`
}

type UserCom struct {
	UserName   string
	UserId     uint
	UserAvatar string
}
type Like struct {
	UserID  uint `json:"userId"`
	VideoID uint `json:"vedioId"`
}

type Favorite struct {
	UserID  uint `json:"userId"`
	VideoID uint `json:"vedioId"`
}

type SetVedioRequest struct {
	Title        string `json:"title" binding:"required"`
	Url          string `json:"url" binding:"required"`
	Introduction string `json:"introduction"`
	Img          string `json:"img" binding:"required"`
}

type Resolution struct {
	Width  int
	Height int
}

type CommentReplay struct {
	UserName  string    `json:"username"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdtime"`
}
