package controllers

import (
	"backend/models"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var resolutions = map[string]models.Resolution{
	"360p": {640, 360},
	"720p": {1280, 720},
}

// 创建视频
func (ac *AuthController) SetVideo(c *gin.Context) {
	contentHeader := c.GetHeader("Content-Type")

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	if contentHeader != "multipart/form-data" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "文件格式不正确"})
		c.Abort()
		return
	}
	var SetVedioRequest models.SetVedioRequest
	var err error
	if err = c.ShouldBindJSON(&SetVedioRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}
	var lastID int
	query := "SELECT id FROM vedio ORDER BY id DESC LIMIT 1"
	err = ac.DB.QueryRow(query).Scan(&lastID)
	if err != nil {
		log.Fatal(err)
	}
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上传失败，请检查文件参数"})
		return
	}

	fileHandle, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "打开文件失败"})
		return
	}
	defer fileHandle.Close()

	// 检查原始视频分辨率
	osFile, ok := fileHandle.(*os.File)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件处理失败"})
		return
	}

	originalRes, err := resolutionUpload(osFile, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("获取视频分辨率失败：%v", err)})
		return
	}

	// 确定需要转码的分辨率
	var targetResolutions []models.Resolution
	for _, res := range resolutions {
		if originalRes.Width > res.Width || originalRes.Height > res.Height {
			targetResolutions = append(targetResolutions, res)
		}
	}

	// 保存原始视频到临时目录
	srcPath := fmt.Sprintf("static/videoSave/origin/%d", lastID+1)
	if err := c.SaveUploadedFile(file, srcPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 转码并保存不同分辨率的视频
	for _, res := range targetResolutions {
		outPath := fmt.Sprintf("static/videoSave/tmp/%d_%d", lastID, res.Height)
		cmd := exec.Command("ffmpeg",
			"-i", srcPath,
			"-vf", fmt.Sprintf("scale=%d:%d", res.Width, res.Height),
			"-c:v", "libx264",
			"-crf", "23",
			"-c:a", "aac",
			outPath,
		)
		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("视频转码失败：%v", err)})
			return
		}
	}

	_, err = ac.DB.Exec("INSERT INTO vedio(id, title, url, update_time, introduction, likes, favorites, view, createID) VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
		lastID+1, SetVedioRequest.Title, SetVedioRequest.Url, time.Now(), SetVedioRequest.Introduction, 0, 0, userId)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "视频上传成功"})
}

// 获取视频
func (ac *AuthController) GetVideo(c *gin.Context) {
	videoID := c.Param("id")
	var title, filename, introduction, userName, avatar string
	var update_time time.Time
	var Id, likes, favorites, view, userId int
	err := ac.DB.QueryRow(`
	SELECT v.id, v.title, v.url, v.update_time, v.Introduction, v.likes, v.favorites, v.view, v.userName, v.userId, u.avatar 
	FROM videos v
	LEFT JOIN users u ON v.userId = u.id WHERE v.id=?`, videoID).Scan(
		&Id, &title, &filename, &update_time, &introduction, &likes, &favorites, &view, &userName, &userId, &avatar)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "视频不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}
	formattedUpdateTime := update_time.In(time.Local).Format("2006-01-02 15:04:05")
	fmt.Sprintln("likes:%d", likes)
	// 构造前端可访问的URL
	videoURL := "/static/videoSave/" + filename

	c.JSON(http.StatusOK, gin.H{"videoId": Id, "title": title, "url": videoURL, "update_time": formattedUpdateTime,
		"introduction": introduction, "likes": likes, "favorites": favorites, "view": view,
		"userName": userName, "userId": userId, "avatar": avatar})
}

// 点赞视频
func (ac *AuthController) LikeVideo(c *gin.Context) {
	videoID := c.Param("id")
	type LikeRequest struct {
		UserId int `json:"userId"`
	}
	var req LikeRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}
	userId := req.UserId // 从请求中获取userId
	// 使用INSERT ... ON DUPLICATE KEY UPDATE语法
	// 如果用户未点赞过该视频，则插入新记录，likes=1
	// 如果用户已点赞过该视频，则将likes取反(1变0或0变1)
	_, err = ac.DB.Exec(`
		INSERT INTO likes (userId, videoId, likes) 
		VALUES (?, ?, 1) 
		ON DUPLICATE KEY UPDATE 
		likes = IF(likes = 1, 0, 1)
	`, userId, videoID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "插入失败",
		})
		return
	}
	var like int
	err = ac.DB.QueryRow("SELECT likes FROM likes WHERE userId = ? && videoId = ?", userId, videoID).Scan(&like)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "操作成功", "islike": like})
}

// 收藏视频
func (ac *AuthController) FavoriteVideo(c *gin.Context) {
	videoID := c.Param("id")
	type LikeRequest struct {
		UserId int `json:"userId"`
	}
	var req LikeRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}
	userId := req.UserId
	_, err = ac.DB.Exec(`
		INSERT INTO favorites (userId, videoId, favorites) 
		VALUES (?, ?, 1) 
		ON DUPLICATE KEY UPDATE 
		favorites = IF(favorites = 1, 0, 1)
	`, userId, videoID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "插入失败",
		})
		return
	}
	var favorites int
	err = ac.DB.QueryRow("SELECT favorites FROM favorites WHERE userId = ? && videoId = ?", userId, videoID).Scan(&favorites)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "收藏成功", "isfavorites": favorites})
}

func (ac *AuthController) GetComments(c *gin.Context) {
	videoID := c.Param("id")

	videoId, err := strconv.ParseUint(videoID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var comments []models.Comment
	var parentUserID *uint
	var parentUserName *string
	var ParentUserAvatar *string
	// 使用原生 SQL 查询，并关联 ParentUser 信息
	sql := `
    SELECT c.*,
           pc.userId as parent_user_id,
           pc.userName as parent_user_name,
		   pc.userAvatar as parent_user_avatar

    FROM comments c
    LEFT JOIN comments pc ON c.parentId = pc.id
    WHERE c.videoId = ?
    ORDER BY c.rootId ASC, c.createdAt ASC
	`

	rows, err := ac.DB.Query(sql, videoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.UserId,
			&comment.UserName,
			&comment.UserAvatar,
			&comment.VideoID,
			&comment.CreatedAt,
			&comment.ParentId,
			&comment.RootId,

			&parentUserID,
			&parentUserName,
			&ParentUserAvatar,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse comments"})
			return
		}

		if parentUserID != nil && parentUserName != nil {
			comment.ParentUser = &models.UserCom{
				UserId:     *parentUserID,
				UserName:   *parentUserName,
				UserAvatar: *ParentUserAvatar,
			}
		} else {
			comment.ParentUser = nil
		}

		comments = append(comments, comment)
	}
	c.JSON(http.StatusOK, comments)
}

// 添加评论
func (ac *AuthController) AddComment(c *gin.Context) {
	var input models.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构造评论对象
	comment := models.Comment{
		VideoID:    input.VideoID,
		UserId:     input.UserID,
		UserName:   input.UserName,
		UserAvatar: input.UserAvatar,
		Content:    input.Content,
		ParentId:   input.ParentID,
		CreatedAt:  time.Now(),
	}

	// 开始事务
	tx, err := ac.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法开启事务"})
		return
	}

	// 插入评论（不指定 id，由数据库自动生成）
	result, err := tx.Exec(
		"INSERT INTO comments (content, userId, userName, userAvatar, videoId, createdAt, parentId, rootId) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		comment.Content, comment.UserId, comment.UserName, comment.UserAvatar, comment.VideoID,
		comment.CreatedAt, comment.ParentId, 0, // rootId 初始为 0
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "插入评论失败"})
		return
	}

	// 获取自增生成的 id
	commentID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论 ID 失败"})
		return
	}

	// 如果是主楼评论，更新 rootId 为自身 id
	if input.ParentID == nil {
		_, err = tx.Exec("UPDATE comments SET rootId = ? WHERE id = ?", commentID, commentID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "设置 rootId 失败"})
			return
		}
	} else {
		// 如果是子评论，查询父评论的 rootId 并设置
		var parentRootId uint
		err = tx.QueryRow("SELECT rootId FROM comments WHERE id = ?", *input.ParentID).Scan(&parentRootId)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询父评论 rootId 失败"})
			return
		}
		comment.RootId = parentRootId
		_, err = tx.Exec("UPDATE comments SET rootId = ? WHERE id = ?", parentRootId, commentID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新子评论 rootId 失败"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "事务提交失败"})
		return
	}

	// 返回创建的评论
	comment.ID = uint(commentID)
	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

// 查询用户是否点赞了视频
func (ac *AuthController) IsLike(c *gin.Context) {
	videoID := c.Param("id")
	type LikeRequest struct {
		UserId int `json:"userId"`
	}
	var req LikeRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}
	userId := req.UserId
	query := true
	var likes int

	err = ac.DB.QueryRow("SELECT likes FROM likes WHERE userId = ? AND videoId = ?", req.UserId, videoID).Scan(&likes)
	if err == sql.ErrNoRows {
		likes = 0 // 没有记录则默认未点赞
		query = false
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}
	fmt.Sprintln(likes)
	c.JSON(http.StatusOK, gin.H{"likes": likes, "query": query, "userId": userId})

}

func (ac *AuthController) IsFavorite(c *gin.Context) {
	videoID := c.Param("id")
	type FavoriteRequest struct {
		UserId int `json:"userId"`
	}
	var req FavoriteRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}
	userId := req.UserId

	var favorites int
	err = ac.DB.QueryRow("SELECT favorites FROM favorites WHERE userId = ? AND videoId = ?", userId, videoID).Scan(&favorites)
	if err == sql.ErrNoRows {
		favorites = 0
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func (ac *AuthController) View(c *gin.Context) {
	videoID := c.Param("id")
	var err error
	_, err = ac.DB.Exec("UPDATE videos SET view = view + 1 WHERE id = ?", videoID)
	if err != nil {
		c.JSON(500, gin.H{"message": "更新播放量失败"})
		return
	}
	c.JSON(200, gin.H{"message": "更新成功"})

}

func resolutionUpload(file *os.File, c *gin.Context) (models.Resolution, error) {
	tempFile, err := os.CreateTemp("", "video_*.mp4")
	if err != nil {
		return models.Resolution{}, fmt.Errorf("创建临时文件失败：%v", err)
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	if _, err = io.Copy(tempFile, file); err != nil {
		return models.Resolution{}, fmt.Errorf("写入临时文件失败：%v", err)
	}

	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=p=0",
		tempFile.Name(),
	)
	output, err := cmd.Output()
	if err != nil {
		return models.Resolution{}, fmt.Errorf("解析视频信息失败：%v", err)
	}

	parts := bytes.Split(bytes.TrimSpace(output), []byte(","))
	if len(parts) != 2 {
		return models.Resolution{}, errors.New("无效的分辨率数据")
	}

	width, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return models.Resolution{}, fmt.Errorf("解析宽度失败：%v", err)
	}

	height, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return models.Resolution{}, fmt.Errorf("解析高度失败：%v", err)
	}

	return models.Resolution{Width: width, Height: height}, nil
}
