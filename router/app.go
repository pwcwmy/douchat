package router

import (
	"douchat/service"

	"github.com/gin-gonic/gin"
	"douchat/docs"
   	swaggerfiles "github.com/swaggo/files"
   	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine{
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/index", service.GetIndex)

	r.GET("/user/getUserList", service.GetUserList)
	r.GET("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd) // 登录

	r.GET("/user/sendMsg", service.SendMsg) // 发送消息	

	r.GET("/user/sendUserMsg", service.SendUserMsg) // 发送消息	

	return r
}