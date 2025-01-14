package service

import (
	"douchat/models"
	"douchat/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "获取用户列表成功",
		"data":    data,
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
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名已经注册",
			"data":    user,
		})
		return
	}
	password := c.Query("password")
	repassword := c.Query("repassword")
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "两次输入密码不一致",
			"data":    user,
		})
		return
	}
	// 存入密码时用pwd+salt随机数 Md5加密
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	// user.Password = password
	user.Password = utils.MakePassWord(password, salt)
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "创建用户成功",
		"data":    user,
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
		"code":    0,
		"message": "删除用户成功",
		"data":    user,
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
			"code":    -1,
			"message": "修改用户失败, 电话或邮箱校验不合法",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "修改用户成功",
			"data":    user,
		})
	}
}

// FindUserByNameAndPwd
// @Tags 用户模块
// @Summary 登录
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	// data := make([] *models.UserBasic, 10)
	// data = models.GetUserList()
	data := models.UserBasic{}
	name := c.Query("name")
	// 1. 判断是否存在
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "该用户不存在, 显示的是一个空的UserBasic",
			"data":    data,
		})
		return
	}

	password := c.Query("password")
	// 2. 校验密码
	flag := utils.ValidPassWord(password, user.Salt, user.Password)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	// 查询数据库时所用的密码应该为加密处理后的密码
	data = models.FindUserByNameAndPwd(name, user.Password)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    data,
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		tm := time.Now().Format("2015-01-12 15:06:55")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}
