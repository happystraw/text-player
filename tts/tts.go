package tts

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// XunfeiTts Xunfei tts struct
type XunfeiTts struct {
	Host   string
	AppID  string
	APIKey string
}

// BaseTts interface
type BaseTts interface {
	Create(msg string, params map[string]string) ([]byte, error)
}

// New a xunfei tts instance
func New(host, appid, apikey string) *XunfeiTts {
	return &XunfeiTts{
		Host:   host,
		AppID:  appid,
		APIKey: apikey,
	}
}

// Create 请求讯飞接口生层语音合成文件
func (tts *XunfeiTts) Create(msg string, params map[string]string) ([]byte, error) {
	curtime := strconv.FormatInt(time.Now().Unix(), 10)

	param, _ := json.Marshal(params)
	base64Param := base64.StdEncoding.EncodeToString(param)

	w := md5.New()
	io.WriteString(w, tts.APIKey+curtime+base64Param)
	checksum := fmt.Sprintf("%x", w.Sum(nil))

	var data = url.Values{}
	data.Add("text", msg)

	reqBody := data.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("POST", tts.Host, strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("X-CurTime", curtime)
	req.Header.Set("X-Appid", tts.AppID)
	req.Header.Set("X-Param", base64Param)
	req.Header.Set("X-CheckSum", checksum)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if respOk, ok := resp.Header["Content-Type"]; !ok || len(respOk) < 1 || respOk[0] == "text/plain" {
		return []byte{}, errors.New("Error: xunfei tts api create fail: " + string(respBody))
	}

	return respBody, nil
}
