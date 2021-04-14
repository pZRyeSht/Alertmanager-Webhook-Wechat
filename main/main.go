package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"project/prometheus-alert/model"
	"project/prometheus-alert/notifier"
)

var (
	help        bool
	RobotKey string
)

func init() {
	flag.BoolVar(&help, "help", false, "help")
	flag.StringVar(&RobotKey, "RobotKey", "", "global wechat robot webhook, you can overwrite by alert rule with annotations wechatRobot")
}

func main() {

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	router := gin.Default()
	router.POST("/webhook", func(c *gin.Context) {
		var notification model.Notification

		err := c.BindJSON(&notification)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		RobotKey := c.DefaultQuery("key", RobotKey)

		err = notifier.Send(notification, RobotKey)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}

		c.JSON(http.StatusOK, gin.H{"message": "send to wechat robot successful!"})

	})
	_ = router.Run(":6666")
}
