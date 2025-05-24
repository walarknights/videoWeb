package controllers

import (
	"backend/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings" // 引入 strings 包
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// --- 统一密钥管理 ---
// 建议从配置或环境变量加载
const jwtSecretKey = "walark"

// --- 核心 Token 验证函数 ---
// parseAndValidateToken 解析并验证 JWT token 字符串
// 返回 claims 或 error
func parseAndValidateToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名算法: %v", token.Header["alg"])
		}
		// 返回用于验证签名的密钥
		return secretKey, nil
	})

	// 处理解析错误 (包括过期等)
	if err != nil {
		log.Printf("token 解析失败: %v", err)   // 记录详细错误
		return nil, fmt.Errorf("无效或已过期的令牌") // 返回通用错误给调用者
	}

	// 验证 token 是否有效且 claims 是 jwt.MapClaims 类型
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 可选：在这里可以添加对必要 claims (如 userId, username) 是否存在的检查
		_, userIdExists := claims["userId"]
		_, usernameExists := claims["username"]
		if !userIdExists || !usernameExists {
			return nil, fmt.Errorf("令牌缺少必要的声明(claims)")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("无效的令牌")
}

// --- AuthController 结构体和 New 函数保持不变 ---
type AuthController struct {
	DB *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{DB: db}
}

// --- Login 和 Setup 函数保持不变 (只展示 Login 示例) ---
func (ac *AuthController) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	var err error
	if err = c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}

	var user models.User
	query := "SELECT id, username, password, followers, followings, avatar, dynamicNum FROM users WHERE username = ?"
	err = ac.DB.QueryRow(query, loginReq.Username).Scan(
		&user.ID, &user.Username, &user.Password,
		&user.Followers, &user.Following, &user.Avatar, &user.DynamicNum,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		} else {
			log.Printf("数据库查询错误: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		}
		return
	}

	corret := checkPassword(user.Password, loginReq.Password)
	if !corret {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	// 生成JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 48).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecretKey)) // 使用常量密钥
	if err != nil {
		log.Printf("无法生成令牌: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		return
	}

	responseUser := models.UserResponse{
		Username:   user.Username,
		Followers:  user.Followers,
		Following:  user.Following,
		Avatar:     user.Avatar,
		DynamicNum: user.DynamicNum,
		UserId:     user.ID,
		IsLogin:    true,
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenString,
		User:  responseUser,
	})
}

// --- 修改后的 VerifyToken 控制器 ---
func (ac *AuthController) VerifyToken(c *gin.Context) {
	var tokenReq models.TokenRequest
	if err := c.ShouldBindJSON(&tokenReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数", "valid": false})
		return
	}

	// 1. 使用核心函数验证 Token 签名和有效期
	claims, err := parseAndValidateToken(tokenReq.Token, []byte(jwtSecretKey))
	if err != nil {
		// 核心函数已记录日志，这里直接返回错误信息
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "valid": false})
		return
	}

	// 2. Token 本身有效，现在从 claims 中获取用户信息并查询数据库确认用户存在
	//    注意：JWT 中的数字通常解析为 float64，需要转换
	userIdClaim, okId := claims["userId"].(float64)

	if !okId {
		log.Printf("VerifyToken: 令牌声明类型错误或缺失: userId=%v (%T), username=%v (%T)", claims["userId"], claims["userId"], claims["username"], claims["username"])
		c.JSON(http.StatusUnauthorized, gin.H{"message": "无效的令牌声明", "valid": false})
		return
	}
	userId := int(userIdClaim) // 转换为 int (或 int64，取决于你的 User ID 类型)

	// 3. 查询数据库 (可选，但 VerifyToken 的目的通常包含此检查)
	var user models.User
	query := "SELECT id, username, password, followers, followings, avatar, dynamicNum FROM users WHERE id = ?"
	err = ac.DB.QueryRow(query, userId).Scan(
		&user.ID, &user.Username, &user.Password,
		&user.Followers, &user.Following, &user.Avatar, &user.DynamicNum,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Token 签名有效，但对应的用户已不存在
			c.JSON(http.StatusUnauthorized, gin.H{"message": "与令牌关联的用户不存在", "valid": false})
		} else {
			log.Printf("VerifyToken 数据库查询错误: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误", "valid": false}) // 服务器错误时，不应标记为 valid
		}
		return
	}

	responseUser := models.UserResponse{
		Username:   user.Username,
		Followers:  user.Followers,
		Following:  user.Following,
		Avatar:     user.Avatar,
		DynamicNum: user.DynamicNum,
		UserId:     user.ID,
		IsLogin:    true,
	}
	// 令牌有效，且用户存在
	c.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenReq.Token,
		User:  responseUser,
	})
}

// --- hashPassword 和 checkPassword 保持不变 ---
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// --- 修改后的 AuthMiddleware ---
// 注意：中间件通常不直接依赖 *sql.DB，如果需要数据库检查，需考虑如何注入 DB 或采用不同策略
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请提供认证令牌"})
			c.Abort()
			return
		}

		// 提取 token (去除 "Bearer " 前缀)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader || tokenString == "" { // 检查前缀是否存在且 token 不为空
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误 (需要 Bearer token)"})
			c.Abort()
			return
		}

		// 使用核心函数验证 Token
		claims, err := parseAndValidateToken(tokenString, []byte(jwtSecretKey))
		if err != nil {
			// 核心函数已记录日志，这里直接返回错误信息
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Token 有效，将用户信息存入上下文
		// 同样注意类型转换和检查
		userIdClaim, okId := claims["userId"].(float64)
		usernameClaim, okUsername := claims["username"].(string)

		if !okId || !okUsername {
			log.Printf("AuthMiddleware: 令牌声明类型错误或缺失: userId=%v (%T), username=%v (%T)", claims["userId"], claims["userId"], claims["username"], claims["username"])
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌声明"})
			c.Abort()
			return
		}

		// 将验证后的信息放入 Context，供后续 Handler 使用
		c.Set("userId", int(userIdClaim)) // 转换为合适的类型
		c.Set("username", usernameClaim)

		// 请求继续向下传递
		c.Next()
	}
}

func (ac *AuthController) Setup(c *gin.Context) {

	var SetupReq models.SetupRequest
	var err error
	if err = c.ShouldBindJSON(&SetupReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求参数"})
		return
	}
	var lastID int
	query := "SELECT id FROM users ORDER BY id DESC LIMIT 1"
	err = ac.DB.QueryRow(query).Scan(&lastID)
	// 注意：这里如果在空表上查询会出错，需要处理 sql.ErrNoRows
	if err != nil && err != sql.ErrNoRows {
		log.Printf("获取 lastID 错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		return
	}
	// 如果是空表，lastID 默认为 0，新 ID 为 1

	stmt, err := ac.DB.Prepare("INSERT INTO users(id, username, password, followers, following, avatar, dynamicNum) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil { // 检查 Prepare 错误
		log.Printf("数据库 Prepare 错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		return
	}
	defer stmt.Close()

	hPassword, err := hashPassword(SetupReq.SetPassword)
	if err != nil { // 检查 hashPassword 错误
		log.Printf("密码哈希错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		return
	}

	newUserID := lastID + 1
	_, err = stmt.Exec(newUserID, SetupReq.SetUsername, hPassword, 0, 0, "https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg", 0)
	if err != nil {
		// 检查可能的数据库错误，例如用户名唯一性冲突
		log.Printf("数据库插入错误: %v", err)
		// 可以根据具体的数据库错误类型返回更友好的信息，如用户名已存在
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注册用户失败"})
		return
	}

	// 生成JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   newUserID, // 使用新用户的 ID
		"username": SetupReq.SetUsername,
		"exp":      time.Now().Add(time.Hour * 48).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecretKey)) // 使用常量密钥
	if err != nil {
		log.Printf("无法生成令牌: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		return
	}

	responseUser := models.UserResponse{
		Username:   SetupReq.SetUsername,
		Followers:  0,
		Following:  0,
		Avatar:     "https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg",
		DynamicNum: 0,
		UserId:     newUserID, // 使用新用户的 ID
		IsLogin:    true,
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenString,
		User:  responseUser,
	})
}
