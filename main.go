package main
import (
	"douchat/router"
	"douchat/utils"
)
func main() {
	utils.InitConfig()
	utils.InitMySQL()
	r := router.Router()
	r.Run(":8081")
}