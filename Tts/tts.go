package Tts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	apiKey    = "k1MZ7P6gB27TeFvn1yU62rGs"
	secretKey = "b9iS1zykWPAl7NklVE15bLE6GGLuCTTf"
	openApi   = "https://openapi.baidu.com/oauth/2.0/token"
	ttsApi    = "http://tsn.baidu.com/text2audio"
)

func Conver(text string) (string, error) {
	accessToken, err := getToken()
	if err != nil {
		return "", err
	}
	//发音人选择, 0为普通女声，1为普通男生，3为情感合成-度逍遥，4为情感合成-度丫丫，默认为普通女声
	per := "0"
	//语速，取值0-15，默认为5中语速
	spd := "5"
	//音调，取值0-15，默认为5中语调
	pit := "5"
	//音量，取值0-9，默认为5中音量
	vol := "6"
	// 下载的文件格式, 3：mp3(default) 4： pcm-16k 5： pcm-8k 6. wav
	aue := "6"
	cuid := "123456GO"
	tex := "tex=" + text + "&lan=zh&ctp=1&cuid=" + cuid + "&tok=" + accessToken + "&per=" + per + "&spd=" + spd + "&pit=" + pit + "&vol=" + vol + "&aue=" + aue
	tts_url := ttsApi + "?" + tex
	fileName, err := httpGetTts(tts_url)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

var TokenRes struct {
	AccessToken   string `json:"access_token"`
	ExpiresIn     int64  `json:"expires_in"`
	RefreshToken  string `json:"refresh_token"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
	AccessTime    int64  `json:"-"`
}

func getToken() (string, error) {
	if len(TokenRes.AccessToken) > 0 && (TokenRes.AccessTime+TokenRes.ExpiresIn) < time.Now().Unix() {
		return TokenRes.AccessToken, nil
	}
	auth_url := openApi + "?grant_type=client_credentials&client_id=" + apiKey + "&client_secret=" + secretKey
	err := httpGet(auth_url, &TokenRes)
	if err != nil {
		return "", err
	}
	TokenRes.AccessTime = time.Now().Unix()
	return TokenRes.AccessToken, nil
}

func httpPost(url string, body []byte, v interface{}) error {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("http.Status: %s", httpResp.Status))
	}
	err = json.NewDecoder(httpResp.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func httpGet(url string, v interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("http.Status: %s", httpResp.Status))
	}
	err = json.NewDecoder(httpResp.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func httpGetTts(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("http.Status: %s", httpResp.Status))
	}

	if head := httpResp.Header.Get("Content-Type"); head == "audio/wav" {
		fileName := randSeq(8) + ".wav"
		f, err := os.Create(fileName)
		if err != nil {
			return "", err
		}
		defer f.Close()
		io.Copy(f, httpResp.Body)
		return fileName, nil
	}

	var Res struct {
		TtsLogid   float64 `json:"tts_logid"`
		ErrDetail  string  `json:"err_detail"`
		ErrMsg     string  `json:"err_msg"`
		ErrNo      int     `json:"err_no"`
		ErrSubcode int     `json:"err_subcode"`
	}
	err = json.NewDecoder(httpResp.Body).Decode(&Res)
	if err != nil {
		return "", err
	}
	return "", errors.New(Res.ErrMsg)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
