package router

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoute() *gin.Engine {
	r := gin.Default()
	// 配置静态文件目录
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I

	apiRouter.POST("/favorite/action/", middleware.JwtAuth(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.JwtAuth(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.JwtAuth(), controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.JwtAuth(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.JwtAuth(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.JwtAuth(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.JwtAuth(), controller.FollowerList)
	// apiRouter.GET("/relation/friend/list/", controller.FriendList)
	// apiRouter.GET("/message/chat/", controller.MessageChat)
	// apiRouter.POST("/message/action/", controller.MessageAction)
	return r
}
