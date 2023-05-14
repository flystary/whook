package goclient

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"web/modules/gojson"
	"web/modules/hash"


	"github.com/gorilla/websocket"
)

type WebHookClient struct {
	token string
	conn *websocket.Conn

}

func (p *WebHookClient)Connection(remoteUrl  string,token string) (*WebHookClient ,error){
	client :=  &WebHookClient{

	}

	header := http.Header{}
	header.Add("x-smarthook-token",token)

	c, _, err := websocket.DefaultDialer.Dial(remoteUrl,header)

	if err != nil {
		return client,err
	}

	client.conn = c
	return client,nil
}

func (p *WebHookClient) SetCloseHandler(h func(code int, text string) error) {
	p.conn.SetCloseHandler(h)
}

func (p *WebHookClient)Send(msg []byte) error {

	return p.conn.WriteMessage(websocket.TextMessage, msg)
}

func (p *WebHookClient) SendJSON(v interface{}) error {
	return p.conn.WriteJSON(v)
}

func (p *WebHookClient) Read() ([]byte,error) {
	_, message, err := p.conn.ReadMessage()

	return message,err
}

func (p *WebHookClient) ReadJSON(v interface{}) error {
	return p.conn.ReadJSON(v)
}
func (p *WebHookClient) Close() {
	p.conn.Close()
}

// GetToken 获取与服务器端连接的认证密钥.
func GetToken(remoteUrl,account string,password string) (string,error) {

	t := strconv.Itoa(time.Now().Nanosecond())
	h := sha256.New()
	h.Write([]byte(account + password + t))
	md := h.Sum(nil)
	mdStr := hex.EncodeToString(md)

	response, err := http.Post(remoteUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader("account=" + account +"&password=" + mdStr + "&time=" + t))
	if err != nil {
		return "",err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "",err
	}

	if js := gojson.DeserializeObject(string(body)).GetJsonObject("error_code");js.IsValid() {
		errorCode,err := strconv.Atoi(js.ToString())

		if err != nil {
			return "",err
		}
		if errorCode != 0 {
			message := gojson.DeserializeObject(string(body)).GetJsonObject("message")

			return "",errors.New(message.ToString())
		}
		token := gojson.DeserializeObject(string(body)).GetJsonObject("data")

		return token.ToString(),nil

	}
	return "",errors.New("Data error:" + string(body))

}

// Command 执行命令.
func (p *WebHookClient) Command (host url.URL,account,password ,shell string,channel chan <-[]byte) {

	defer close(channel)


	token,err := GetToken(host.String() +"/token",account,password)

	if err != nil {
		logs.Error("Connection remote server error:", err.Error())

		channel <- []byte("Error: Connection remote server error => " + err.Error())
		return
	}


	u := &url.URL{Scheme: "ws", Host: host.Host , Path: "/socket"}

	client,err := (&WebHookClient{}).Connection(u.String(),token)

	if err != nil {
		logs.Error("Remote server error:", err.Error())

		channel <- []byte("Error:Remote server error => " + err.Error())
		return
	}

	defer client.Close()

	client.SetCloseHandler(func(code int, text string) error {

		return nil
	})

	msgId :=  hash.Md5(shell + time.Now().String())

	command := JsonResult{
		ErrorCode	:   0,
		Message		: "ok",
		Command		: "shell",
		MsgId		:   msgId,
		Data		:    shell,
	}


	err = client.SendJSON(command)

	if err != nil {
		logs.Error("Remote server error:", err.Error())

		channel <- []byte("Error:Remote server error => " + err.Error())
		return
	}

	for {
		var response JsonResult

		err := client.ReadJSON(&response)

		if err != nil {
			logs.Error("Remote server error:", err.Error())

			channel <- []byte("Error:Remote server error => " + err.Error())
			return
		}
		if response.ErrorCode == 0 {
			if response.Command == "end" {
				return
			}
			body := response.Data.(string)

			channel <- []byte(body)
		}
	}
}