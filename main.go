package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"golang.design/x/clipboard"
	"io/ioutil"
	"os/exec"
	"time"
)

type config struct {
	RedisAddr     string
	RedisPassword string
	TopicKey      string
	DeviceToken   string
}

var flag = false

var RedisAddr = ""
var RedisPassword = ""
var DeviceToken = ""
var TopicKey = ""

func PublishRedis(client *redis.Client) {
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		if flag == false {
			clipboardContent := string(data)
			encodeString := base64.StdEncoding.EncodeToString([]byte(clipboardContent))
			jsondata := make(map[string]string)
			jsondata["type"] = "text"
			jsondata["msg"] = encodeString
			jsondata["uuid"] = DeviceToken
			dataType, _ := json.Marshal(jsondata)
			dataString := string(dataType)
			client.Publish(TopicKey, dataString)
		}

	}
}

func SubscribeRedis(client *redis.Client) {
	pubsub := client.Subscribe(TopicKey)
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		message := gjson.Get(msg.Payload, "msg").String()
		uuid := gjson.Get(msg.Payload, "uuid").String()
		if uuid != DeviceToken {
			decodeString, _ := base64.StdEncoding.DecodeString(message)
			log.Infof("向剪切板写入信息: %s\n", decodeString)
			clipboard.Write(clipboard.FmtText, decodeString)
			notifyCommand := fmt.Sprintf("display notification \"%s\" with title \"SyncClipBoard同步消息\"", decodeString)
			command := exec.Command("osascript", "-e", notifyCommand)
			err := command.Run()
			if err != nil {
				fmt.Println(err.Error())
			}
			flag = true
			time.Sleep(time.Second * 1)
			flag = false
		}
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	var Config config
	bytes, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(bytes, &Config)

	if err != nil {
		panic(err)
	}
	RedisAddr = Config.RedisAddr
	RedisPassword = Config.RedisPassword
	DeviceToken = Config.DeviceToken
	TopicKey = Config.TopicKey

	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword,
	})

	go SubscribeRedis(client)
	PublishRedis(client)

}
