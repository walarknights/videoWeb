package controllers

import (
	"backend/models"

	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"   // 若需显式使用 *mysql.MySQLError 类型
	_ "github.com/go-sql-driver/mysql" // 注册 MySQL 驱动
)

func (ac *AuthController) GetPersonalInfo(c *gin.Context) {
	userID := c.Param("userId")
	userId, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userInfo models.User
	query := `SELECT * FROM users WHERE id = ?`
	err = ac.DB.QueryRow(query, userId).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password,
		&userInfo.Followers, &userInfo.Following, &userInfo.Avatar, &userInfo.DynamicNum)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})

		} else {
			c.JSON(500, gin.H{"error": err.Error()})

		}
		return
	}

	c.JSON(200, userInfo)
}

func (ac *AuthController) AddFocus(c *gin.Context) {
	var focus struct {
		FocusId   int
		FocusedId int
	}
	var isFocus int
	err := c.ShouldBindJSON(&focus)
	if err != nil {
		c.JSON(400, gin.H{"message": "错误的数据格式"})
		return
	}
	if focus.FocusId == focus.FocusedId {
		isFocus = 0
		c.JSON(200, gin.H{"isFocus": isFocus})
		return
	}
	_, err = ac.DB.Exec("INSERT INTO focus(focusId, focusedId) VALUES (?, ?)", focus.FocusId, focus.FocusedId)

	if err != nil {
		// 判断是否是唯一键冲突
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			isFocus = 1
			c.JSON(200, gin.H{"isFocus": isFocus})
			return
		} else {
			c.JSON(500, gin.H{"error": "插入数据失败"})
			return
		}
	}
	isFocus = 2
	c.JSON(200, gin.H{
		"message": "插入成功",
		"isFocus": isFocus,
	})

}

func (ac *AuthController) GetFoucsList(c *gin.Context) {
	userID := c.Param("userId")
	userId, err := strconv.ParseUint(userID, 10, 32)
	var focusUserId []int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	rows, err := ac.DB.Query(`SELECT focusedId FROM focus WHERE focusId = ?`, userId)
	if err != nil {
		if err == sql.ErrNoRows {

			c.Status(204)
		} else {
			log.Printf("数据库查询错误: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		}
		return
	}
	defer rows.Close()
	for rows.Next() {
		var Id int
		err := rows.Scan(
			&Id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse focus"})
			return
		}
		focusUserId = append(focusUserId, Id)
	}
	type userFocusInfo struct {
		UserName string `json:"userName"`
		Avatar   string `json:"avatar"`
		UserId   int    `json:"userId"`
	}
	var response []userFocusInfo

	// 1. 动态生成占位符
	placeholders := make([]string, len(focusUserId))
	args := make([]interface{}, len(focusUserId))
	for i, id := range focusUserId {
		placeholders[i] = "?"
		args[i] = id

	}

	// 2. 拼接 SQL 查询语句
	query := fmt.Sprintf("SELECT username, avatar, id FROM users WHERE id IN (%s)", strings.Join(placeholders, ","))

	// 3. 执行查询
	rows, err = ac.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	defer rows.Close()
	for rows.Next() {

		var res userFocusInfo
		err := rows.Scan(
			&res.UserName,
			&res.Avatar,
			&res.UserId,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse focus"})
			return
		}

		response = append(response, res)
	}
	c.JSON(200, response)
}
