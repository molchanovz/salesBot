package gigaChat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type CompletionData struct {
	Model             string `json:"model"`
	Messages          []*Msg `json:"messages"`
	Stream            bool   `json:"stream"`
	RepetitionPenalty int64  `json:"repetition_penalty"`
}

type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	msgUrl   = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	tokenUrl = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	authKey  = "Mzg3OGYyMTQtMzM3NS00MTY0LWFkMTEtYzAxZmI2NjNkOWY5OmFmNzY1NjU4LWE5ODUtNDIwMC05MjNhLTdjYmUzODAzNmVmNw=="
	clientId = "3878f214-3375-4164-ad11-c01fb663d9f9"
)

func SendRequest(content string) (Response, error) {
	payload := CompletionData{
		Model: "GigaChat",
		Messages: []*Msg{
			&Msg{
				Role:    "user",
				Content: content,
			},
		},
		Stream:            false,
		RepetitionPenalty: 1,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", msgUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+getToken())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var someResponse Response
	err = json.Unmarshal(data, &someResponse)
	if err != nil {
		return someResponse, err
	}

	return someResponse, nil
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func getToken() string {

	data := make(map[string]string)
	data["scope"] = "GIGACHAT_API_PERS"

	// Преобразуем данные в формат x-www-form-urlencoded
	values := make(url.Values)
	for k, v := range data {
		values.Set(k, v)
	}
	postData := values.Encode()

	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBufferString(postData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+authKey)
	req.Header.Add("RqUID", clientId)

	// Устанавливаем индикатор insecure
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var t Token

	err = json.Unmarshal(res, &t)
	if err != nil {
		println(err)
		return ""
	}

	//println(t.AccessToken, t.ExpiresAt)

	return t.AccessToken
}

type Response struct {
	Choices []struct {
		Message struct {
			Role           string `json:"role"`
			Content        string `json:"content"`
			DataForContext []struct {
			} `json:"data_for_context"`
		} `json:"message"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Object string `json:"object"`
}
