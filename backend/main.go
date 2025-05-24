package main

import (
	"backend/controllers"
	"backend/middleware"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// --- 数据库连接 ---
	// TODO: 从配置或环境变量读取DSN
	dsn := "root:Fast19930108@tcp(localhost:3306)/vedioWeb?parseTime=true" // 建议添加 parseTime=true
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("打开数据库错误: %v\n", err)
	}
	defer db.Close()

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("连接数据库错误: %v\n", err)
	}
	fmt.Println("成功连接到数据库!")

	authController := controllers.NewAuthController(db) // 使用构造函数

	r := gin.Default()
	r.Static("/static", "E:\\goCode\\backend\\static")

	// 中间件
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		// 使用 authController 的 Login 方法作为处理函数
		api.POST("/login", authController.Login)
		api.POST("/setup", authController.Setup)
		api.POST("verify-token", authController.VerifyToken)
	}
	videoGroup := r.Group("/videos")
	{
		videoGroup.GET("/:id", authController.GetVideo)                // 获取视频
		videoGroup.POST("/:id/like", authController.LikeVideo)         // 点赞视频
		videoGroup.POST("/:id/favorite", authController.FavoriteVideo) // 收藏视频
		videoGroup.GET("/:id/comments", authController.GetComments)    // 获取评论
		videoGroup.POST("/:id/comments", authController.AddComment)    // 添加评论
		videoGroup.POST("/:id/video", authController.SetVideo)         // 发布视频
		videoGroup.POST("/:id/islike", authController.IsLike)          //检查是否点赞
		videoGroup.POST("/:id/isfavorite", authController.IsFavorite)  //检查是否收藏
		videoGroup.POST("/:id/view", authController.View)              //增加播放量
	}
	personalGroup := r.Group("/personal")
	{
		personalGroup.GET("/:userId/userInfo", authController.GetPersonalInfo) //获取基本的用户信息
		personalGroup.POST("/:userId/addFocus", authController.AddFocus)       //增加关注
		personalGroup.GET("/:userId/focusList", authController.GetFoucsList)   //获取关注列表
	}
	// 启动服务器
	fmt.Println("服务器运行在 :8080")
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
