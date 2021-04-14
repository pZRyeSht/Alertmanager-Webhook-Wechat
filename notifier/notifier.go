package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"project/prometheus-alert/model"
	"project/prometheus-alert/transformer"
)

// Send send markdown message to wechat robot
func Send(notification model.Notification, defaultRobot string) (err error) {

	markdown, robotURL, err := transformer.TransformToMarkdown(notification)
	if err != nil {
		log.Println("level=error, TransformToMarkdown Error", err)
		return
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		log.Println("level=error, json.Marshal Error", err)
		return
	}

	var wechatRobotURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + defaultRobot
	if robotURL != "" {
		wechatRobotURL = robotURL
	}

	req, err := http.NewRequest(
		"POST",
		wechatRobotURL,
		bytes.NewBuffer(data))
	if err != nil {
		log.Println("level=error, Post Request Error", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("level=error, Set Client Error", err)
		return
	}

	defer resp.Body.Close()
	fmt.Println("level=info, response Status:", resp.Status)
	fmt.Println("level=info, response Headers:", resp.Header)

	return
}
