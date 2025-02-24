package gpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"wechatbot/config"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// const BASEURL = "https://api.openai.com/v1/"

// ChatGPTResponseBody 响应体
type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ChoiceItem struct {
}

// ChatGPTRequestBody 请求体
type ChatGPTRequestBody struct {
	Model    string              `json:"model"`
	Messages []ChatGPTChatFormat `json:"messages"`
	// Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

type ChatGPTChatFormat struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completions gtp文本模型回复
// curl https://api.openai.com/v1/chat/completions \
//   -H 'Content-Type: application/json' \
//   -H 'Authorization: Bearer YOUR_API_KEY' \
//   -d '{
//   "model": "gpt-3.5-turbo",
//   "messages": [{"role": "user", "content": "Hello!"}]
// }'

// 老版本，但是openai开了全局代理之后无法正常使用
func Completions(msg string) (string, error) {

	chatformat := make([]ChatGPTChatFormat, 0)
	chatformat = append(chatformat, ChatGPTChatFormat{
		Role:    "user",
		Content: msg,
	})

	requestBody := ChatGPTRequestBody{
		// Model:  "text-davinci-003",
		// Prompt: msg,
		Model:            "gpt-3.5-turbo",
		Messages:         chatformat,
		MaxTokens:        2048,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gpt json string : %v", string(requestData))
	req, err := http.NewRequest("POST", (config.Config.EndPoint+"/v1/chat/completions"), bytes.NewBuffer(requestData))
	if err != nil {
		return "http.NewRequest ", err
	}

	apiKey := config.Config.ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "client.Do ", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "ioutil.ReadAll ", err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			var slice map[string]interface{}
			
			slice = v["message"].(map[string]interface{})
			reply = slice["content"].(string)
			break
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}

// 船新版本，使用go的exec库调用python的openai库就可以跑
func CompletionsNew(msg string) (string, error) {

	log.Printf("start command to call python script")
	reply, err := exec.Command("python", "get_gpt_reply.py", "-apikey", config.Config.ApiKey, "-msg", msg).Output()
	if err != nil {
		return "exec.Command", err
	}
	reply, err = simplifiedchinese.GBK.NewDecoder().Bytes(reply)
	if err != nil {
		return "simplifiedchinese.GBK.NewDecoder()", err
	}
	result := strings.TrimSpace(string(reply))
	log.Printf("gpt response text: %s \n", result)
	return result, err
}
