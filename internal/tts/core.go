package tts

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/happystraw/text-player/internal/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type StatusFrame = int

const (
	StatusFirstFrame StatusFrame = iota
	StatusContinueFrame
	StatusLastFrame
)

type Tts interface {
	Create(msg string) ([]byte, error)
}

type XunFeiTts struct {
	cfg *config.Tts
}

type Response struct {
	Sid     string       `json:"sid"`
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
}

type ResponseData struct {
	Audio  string `json:"audio"`
	Ced    string `json:"ced"`
	Status int    `json:"status"`
}

func New(cfg *config.Tts) Tts {
	return &XunFeiTts{cfg: cfg}
}

// Create raw audio data from tts server
func (t *XunFeiTts) Create(msg string) ([]byte, error) {
	conn, err := dial(t.cfg.Host, t.cfg.ApiKey, t.cfg.ApiSecret)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	params := map[string]interface{}{
		"common": map[string]string{
			"app_id": t.cfg.AppId,
		},
		"business": t.cfg.Params,
		"data": map[string]interface{}{
			"status": 2,
			"text":   base64.StdEncoding.EncodeToString([]byte(msg)),
		},
	}

	err = conn.WriteJSON(params)
	if err != nil {
		return nil, err
	}

	var data []byte
	for {
		tmp, ok, err := fetch(conn)
		if err != nil {
			return nil, err
		}
		data = append(data, tmp...)
		if ok {
			break
		}
	}

	return data, nil
}

func fetch(conn *websocket.Conn) ([]byte, bool, error) {
	resp := Response{}
	err := conn.ReadJSON(&resp)
	if err != nil {
		return nil, false, err
	}

	if resp.Code != 0 {
		return nil, false, fmt.Errorf("error: api request %s", resp.Message)
	}

	if len(resp.Data.Audio) < 1 {
		return nil, false, fmt.Errorf("error: audio empty")
	}

	data, err := base64.StdEncoding.DecodeString(resp.Data.Audio)
	if err != nil {
		return nil, false, err
	}

	return data, resp.Data.Status == StatusLastFrame, nil
}

// dial connect tts websocket server
func dial(host, apiKey, apiSecret string) (*websocket.Conn, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(assembleAuthUrl(host, apiKey, apiSecret), nil)
	if err != nil {
		return conn, err
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("handshake failed:message=%s,httpCode=%d", string(b), resp.StatusCode)
	}

	return conn, nil
}

func assembleAuthUrl(host string, apiKey, apiSecret string) string {
	ul, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
	}
	date := time.Now().UTC().Format(time.RFC1123)
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	sign := strings.Join(signString, "\n")
	sha := HmacWithShaToBase64(sign, apiSecret)
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	callurl := host + "?" + v.Encode()

	return callurl
}

func HmacWithShaToBase64(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(encodeData)
}
