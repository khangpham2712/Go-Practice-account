package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	acc "test/account"
	"test/config"
)

func main() {
	conf, err := config.ReadFromConfigFile(".")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	context.Background()
	router := gin.Default()
	setEngineGroupAccount(router, conf)
	if err = router.Run(":" + conf.Port); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func setEngineGroupAccount(router *gin.Engine, config config.Config) {
	accountAPI := router.Group("/account")
	accountAPI.GET("/view", acc.GetAllAccounts(config))
	accountAPI.POST("/create", acc.CreateAccount(config))
	accountAPI.POST("/delete/:name", acc.DeleteAccount(config))
}
