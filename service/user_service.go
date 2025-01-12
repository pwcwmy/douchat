package service

import (
	"douchat/models"
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([] *models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

// CreateUser
// @Tags 用户模块
// @Summary 新增用户
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次输入密码不一致",
		})
		return
	}
	user.Password = password
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"message": "创建用户成功",
	})
}

// DeleteUser
// @Tags 用户模块
// @Summary 删除用户
// @param id query string false "id"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Tags 用户模块
// @Summary 修改用户
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "电话"
// @param email formData string false "邮箱"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	// 校验
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println("修改用户失败, 电话或邮箱校验不合法", err)
		c.JSON(200, gin.H{
			"message": "修改用户失败, 电话或邮箱校验不合法",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
		"message": "修改用户成功",
		})
	}
	
}
