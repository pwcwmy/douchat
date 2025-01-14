package main
import (
	"douchat/router"
	"douchat/utils"
)
func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8082")
}