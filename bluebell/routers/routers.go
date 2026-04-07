package routers

import (
	"Bluebell/controller"
	"Bluebell/logger"
	"Bluebell/middlewares"
	"net/http"

	//"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))

	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)

	r.GET("/debug/pool-stats", controller.PoolStats)

	v1.Use(middlewares.JWTAuthMiddleware()) //应用jwt认证中间件

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatPostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		//根据时间或分数获取帖子列表
		v1.GET("/posts2", controller.GetPostListHandler2)
		//投票
		v1.POST("/vote", controller.PostVoteController)
	}

	// r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	// 	//如果是登录用户，判断请求头是否有，有效JWT
	// 	isLogin := true
	// 	c.Request.Header.Get("Authorizatior")
	// 	if isLogin {
	// 		//如果是登录用户
	// 		c.String(http.StatusOK, "pong")
	// 	} else {
	// 		//否则返回请登录
	// 		c.String(http.StatusOK, "请登录")
	// 	}

	// })
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
