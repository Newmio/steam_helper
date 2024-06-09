package hplib

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type ProxyConfig struct {
	Login string
	Pass  string
	IP    string
	Port  string
}

type Param struct {
	Method          string
	Url             string
	Body            []byte
	Headers         map[string]string
	Cookies         []*http.Cookie
	ResetRequest    bool
	SecondsForReset int
	CreateAllLog    bool
	OnlyErrorLog    bool
}

type ICustomHTTP interface {
	Do(param Param) (customResponse, error)
}

type customHTTP struct {
	client        *http.Client
	proxyUrl      []*url.URL
	mu            *sync.Mutex
	indexProxy    int // индекс прокси
	proxyReqCount int // счетчик запросов с прокси
	maxProxyReq   int // максимальное количество запросов с прокси перед сменой на новое прокси
}

type customResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Cookies   []*http.Cookie
}

func NewHttp(client *http.Client, proxy []ProxyConfig) ICustomHTTP {
	var c customHTTP
	c.client = client
	c.mu = &sync.Mutex{}

	for _, value := range proxy {

		url, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", value.Login, value.Pass, value.IP, value.Port))
		if err != nil {
			panic(err)
		}

		c.proxyUrl = append(c.proxyUrl, url)
	}

	return &c
}

func (c *customHTTP) Do(param Param) (customResponse, error) {
	var resp customResponse

	if c.maxProxyReq == c.proxyReqCount {
		c.updateProxy()
	}

	req, err := http.NewRequest(param.Method, param.Url, bytes.NewBuffer(param.Body))
	if err != nil {
		return resp, Trace(err)
	}

	for key, value := range param.Headers {
		req.Header.Add(key, value)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return resp, Trace(err)
	}
	defer response.Body.Close()
	c.proxyReqCount++

	for key, values := range response.Header {
		resp.Headers[key] = strings.Join(values, ", ")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return resp, Trace(err)
	}
	resp.Body = body
	resp.StatusCode = response.StatusCode
	resp.Cookies = response.Cookies()

	return resp, nil
}

func (c *customHTTP) updateProxy() {
	if len(c.proxyUrl) == 0 {
		return
	}

	c.mu.Lock()

	c.indexProxy++
	c.maxProxyReq = rand.Intn(30-10+1) + 10
	c.proxyReqCount = 0

	if c.client.Transport == nil {
		c.client.Transport = &http.Transport{Proxy: http.ProxyURL(c.proxyUrl[c.indexProxy])}
	} else {
		c.client.Transport.(*http.Transport).Proxy = http.ProxyURL(c.proxyUrl[c.indexProxy])
	}

	c.mu.Unlock()
}
