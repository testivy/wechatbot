package gpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"wechatbot/config"
)

const BASEImgURL = "https://api.openai.com/v1/images/"

type ChatImgResponseBody struct {
	Created uint64              `json:"created"`
	Data    []map[string]string `json:"data"`
}
type ChatImgRequestBody struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

func Generations(msg string) (string, error) {

	requestBody := ChatImgRequestBody{
		Prompt: msg,
		N:      1,
		Size:   "256x256",
	}

	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gpt json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEImgURL+"generations", bytes.NewBuffer(requestData))
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

	gptResponseBody := &ChatImgResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Data) > 0 {
		for _, v := range gptResponseBody.Data {
			reply = v["url"]
			break
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
