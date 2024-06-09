package steam_helper

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
	DoWithAutoProxy(param Param) (customResponse, error)
	DoWithProxy(param Param) (customResponse, error)
	Do(param Param) (customResponse, error)
}

type customHTTP struct {
	clientAutoProxy *http.Client
	clientProxy     *http.Client
	client          *http.Client
	proxyUrl        map[string]*url.URL
	staticProxy     *url.URL
	mu              *sync.Mutex
	proxyReqCount   int // счетчик запросов с прокси
	maxProxyReq     int // максимальное количество запросов с прокси перед сменой на новое прокси
}

type customResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Cookies    []*http.Cookie
}

func NewHttp(client *http.Client, proxy []ProxyConfig) ICustomHTTP {
	var c customHTTP
	c.client = client
	c.mu = &sync.Mutex{}

	if len(proxy) > 0 {
		c.clientAutoProxy = client
		c.clientProxy = client

		for _, value := range proxy {

			url, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", value.Login, value.Pass, value.IP, value.Port))
			if err != nil {
				panic(err)
			}

			c.proxyUrl[value.IP] = url
		}
	}

	return &c
}

// Выполнение запроса с автоматическим выбором прокси
func (c *customHTTP) DoWithAutoProxy(param Param) (customResponse, error) {
	if c.maxProxyReq == c.proxyReqCount {
		c.updateAutoProxy()
	}

	c.proxyReqCount++

	return c.do(param, true, false)
}

// Выполнение запроса без автоматического выбора прокси
func (c *customHTTP) DoWithProxy(param Param) (customResponse, error) {
	return c.do(param, false, false)
}

// Выполнение запроса без прокси
func (c *customHTTP) Do(param Param) (customResponse, error) {
	return c.do(param, false, true)
}

func (c *customHTTP) do(param Param, autoClient, noProxy bool) (customResponse, error) {
	var resp customResponse
	var client *http.Client

	if autoClient {
		client = c.clientAutoProxy
	} else if noProxy {
		client = c.client
	} else {
		client = c.clientProxy
	}

	req, err := http.NewRequest(param.Method, param.Url, bytes.NewBuffer(param.Body))
	if err != nil {
		return resp, Trace(err)
	}

	for key, value := range param.Headers {
		req.Header.Add(key, value)
	}

	response, err := client.Do(req)
	if err != nil {
		return resp, Trace(err)
	}
	defer response.Body.Close()

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

func (c *customHTTP) UpdateProxy(url *url.URL) {
	c.mu.Lock()
	c.staticProxy = url
	c.mu.Unlock()
}

func (c *customHTTP) updateAutoProxy() {
	if len(c.proxyUrl) == 0 {
		return
	}

	c.mu.Lock()

	c.maxProxyReq = rand.Intn(30-10+1) + 10
	c.proxyReqCount = 0

	var keys []string
	for k := range c.proxyUrl {
		keys = append(keys, k)
	}

	if c.client.Transport == nil {
		c.client.Transport = &http.Transport{Proxy: http.ProxyURL(c.proxyUrl[keys[rand.Intn(len(keys))]])}
	} else {
		c.client.Transport.(*http.Transport).Proxy = http.ProxyURL(c.proxyUrl[keys[rand.Intn(len(keys))]])
	}

	c.mu.Unlock()
}

func GetRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.93 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/91.0.864.48",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/10.1.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_5_8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:68.0) Gecko/20100101 Firefox/68.0",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_4_11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_3_9) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:61.0) Gecko/20100101 Firefox/61.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.114 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_2_8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:52.0) Gecko/20100101 Firefox/52.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_1_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.107 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_0_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.66 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 9_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.116 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/25.0.1364.172 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 9_8_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 9_7_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/21.0.1180.83 Safari/537.36",
	"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.18",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.15.7) Presto/2.12.388 Version/12.18",
	"Opera/9.80 (X11; Linux x86_64; U; en) Presto/2.8.131 Version/11.11",
	"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.00",
	"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.5.24 Version/10.54",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_7_5) Presto/2.12.388 Version/12.14",
	"Opera/9.80 (X11; Linux i686; U; en) Presto/2.12.388 Version/12.16",
	"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.17",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_6_8) Presto/2.12.388 Version/12.15",
	"Opera/9.80 (X11; Linux x86_64; U; ru) Presto/2.12.388 Version/12.15",
	"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.12.388 Version/12.16",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_6_8) Presto/2.12.388 Version/12.14",
	"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.12.388 Version/12.15",
	"Opera/9.80 (Windows NT 5.1; U; ru) Presto/2.12.388 Version/12.17",
	"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.10.289 Version/11.64",
	"Opera/9.80 (X11; Linux x86_64; U; fr) Presto/2.10.289 Version/11.64",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_7_5) Presto/2.10.289 Version/11.62",
	"Opera/9.80 (X11; Linux x86_64; U; es-ES) Presto/2.10.289 Version/11.64",
	"Opera/9.80 (Windows NT 5.1; U; it) Presto/2.10.289 Version/11.64",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_6_8) Presto/2.9.181 Version/11.52",
	"Opera/9.80 (Windows NT 6.1; U; de) Presto/2.9.168 Version/11.51",
	"Opera/9.80 (X11; Linux i686; U; pl) Presto/2.9.168 Version/11.50",
	"Opera/9.80 (Windows NT 5.1; U; pt-BR) Presto/2.9.168 Version/11.50",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_5_8) Presto/2.9.168 Version/11.50",
	"Opera/9.80 (Windows NT 6.1; U; fr) Presto/2.8.131 Version/11.11",
	"Opera/9.80 (X11; Linux x86_64; U; en-GB) Presto/2.8.131 Version/11.11",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_6_8) Presto/2.8.131 Version/11.11",
	"Opera/9.80 (Windows NT 5.1; U; de) Presto/2.7.62 Version/11.01",
	"Opera/9.80 (X11; Linux i686; U; ru) Presto/2.7.62 Version/11.01",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_5_8) Presto/2.7.62 Version/11.01",
	"Opera/9.80 (Windows NT 6.1; U; es-ES) Presto/2.7.62 Version/11.00",
	"Opera/9.80 (X11; Linux x86_64; U; en-GB) Presto/2.7.62 Version/11.00",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_4_11) Presto/2.7.62 Version/11.00",
	"Opera/9.80 (Windows NT 5.1; U; en-GB) Presto/2.7.62 Version/11.00",
	"Opera/9.80 (X11; Linux x86_64; U; ru-RU) Presto/2.7.62 Version/11.00",
	"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.7.62 Version/11.00",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_5_8) Presto/2.7.62 Version/10.70",
	"Opera/9.80 (X11; Linux x86_64; U; en) Presto/2.6.30 Version/10.63",
	"Opera/9.80 (Windows NT 6.1; U; es-ES) Presto/2.6.30 Version/10.63",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_5_8) Presto/2.6.30 Version/10.63",
	"Opera/9.80 (X11; Linux x86_64; U; en) Presto/2.6.30 Version/10.61",
	"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.6.30 Version/10.61",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_5_8) Presto/2.6.30 Version/10.61",
	"Opera/9.80 (X11; Linux x86_64; U; en-GB) Presto/2.6.30 Version/10.60",
	"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.6.30 Version/10.60",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_6_8) Presto/2.6.30 Version/10.60",
	"Opera/9.80 (Windows NT 6.1; U; fr) Presto/2.6.30 Version/10.54",
	"Opera/9.80 (X11; Linux x86_64; U; en) Presto/2.6.30 Version/10.54",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10_5_8) Presto/2.6.30 Version/10.54",
	"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.6.30 Version/10.53",
	"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 9; Pixel 3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.93 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 8.0.0; SM-G950F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.101 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 10; SM-N975U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.101 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 13_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 9; Pixel 2 XL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.119 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 12_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 8.1.0; Pixel XL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 7.0; Nexus 5X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 11_4 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 6P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 5.1; Nexus 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.124 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.0 Mobile/14G60 Safari/602.1",
	"Mozilla/5.0 (Linux; Android 4.4.4; Nexus 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.111 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 10_2 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/14C92 Safari/602.1",
	"Mozilla/5.0 (Linux; Android 9; SAMSUNG SM-G960U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/10.1 Chrome/71.0.3578.99 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_3_5 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13G36 Safari/601.1",
	"Mozilla/5.0 (Linux; Android 8.0.0; SAMSUNG SM-G950U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/7.2 Chrome/59.0.3071.125 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 9_3_3 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13G34 Safari/601.1",
	"Mozilla/5.0 (Linux; Android 7.1.2; SM-T580 Build/NJH47F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.158 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H143 Safari/600.1.4",
	"Mozilla/5.0 (Linux; Android 7.0; SAMSUNG SM-G930F) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/5.4 Chrome/51.0.2704.106 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B410 Safari/600.1.4",
	"Mozilla/5.0 (Linux; Android 6.0.1; SAMSUNG SM-T810) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/4.0 Chrome/44.0.2403.133 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 7_1_2 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D257 Safari/9537.53",
	"Mozilla/5.0 (Linux; Android 5.1.1; SAMSUNG SM-T700) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53",
	"Mozilla/5.0 (Linux; Android 4.4.2; SAMSUNG SM-T230) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/2.1 Chrome/34.0.1847.76 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_6 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10B500 Safari/8536.25",
	"Mozilla/5.0 (Linux; Android 4.2.2; SAMSUNG SM-T210) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/1.5 Chrome/28.0.1500.94 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A403 Safari/8536.25",
	"Mozilla/5.0 (Linux; Android 4.1.2; SAMSUNG GT-N8013) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/1.0 Chrome/18.0.1025.166 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 5_1_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9B206 Safari/7534.48.3",
	"Mozilla/5.0 (Linux; Android 4.4.4; Nexus 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.117 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",
	"Mozilla/5.0 (Linux; Android 4.2.2; Nexus 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 4_3_5 like Mac OS X) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8L1 Safari/6533.18.5",
	"Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.138 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 4_3 like Mac OS X) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8F190 Safari/6533.18.5",
	"Mozilla/5.0 (Linux; Android 4.0.4; Nexus S) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.166 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 3_1_3 like Mac OS X) AppleWebKit/528.18 (KHTML, like Gecko) Version/4.0 Mobile/7E18 Safari/528.16",
	"Mozilla/5.0 (Linux; Android 3.2; Nexus 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.59 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 3_2 like Mac OS X) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B367 Safari/531.21.10",
	"Mozilla/5.0 (Linux; Android 2.3.7; Nexus One) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.72 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 2_2_1 like Mac OS X) AppleWebKit/525.20 (KHTML, like Gecko) Version/3.0 Mobile/7F190 Safari/525.20",
	"Mozilla/5.0 (Linux; Android 1.6; HTC Dream) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	"Opera/9.80 (Android 4.2.1; Linux; Opera Mobi/ADR-1302181428) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306260857) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/ADR-1202011742) Presto/2.9.201 Version/11.50",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306272045) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/ADR-1202011742) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.2.2; Linux; Opera Mobi/ADR-1307021151) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.3; Linux; Opera Mobi/ADR-1305181055) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1305250714) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1303201732) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.2.2; Linux; Opera Mobi/ADR-1306271417) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306041641) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 2.3.3; Linux; Opera Mobi/ADR-1306262106) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.2.1; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.3; Linux; Opera Mobi/ADR-1202011737) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1202011737) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.2.2; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1202011739) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.0.3; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.3; Linux; Opera Mobi/ADR-1202011739) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1202011737) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.2.2; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.2.2; Linux; Opera Mobi/ADR-1202011737) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.2.1; Linux; Opera Mobi/ADR-1202011739) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1202011737) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1202011737) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.2.1; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1202011737) Presto/2.10.254 Version/12.00",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.3; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.1; Linux; Opera Mobi/ADR-1306181741) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.0.4; Linux; Opera Mobi/ADR-1306181740) Presto/2.11.355 Version/12.10",
	"Opera/9.80 (Android 4.1.2; Linux; Opera Mobi/ADR-1306181742) Presto/2.11.355 Version/12.10",
}
