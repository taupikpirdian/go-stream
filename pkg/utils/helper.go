package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

func NewHttpClient(timeout time.Duration, proxy string) *http.Client {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if proxy != "" {
		proxyUrl, _ := url.Parse(proxy)
		customTransport.Proxy = http.ProxyURL(proxyUrl)
	}
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport, Timeout: timeout}
	return client
}

func UnixToTime(unix int64) time.Time {
	unixToTime := time.Unix(unix, 0)
	dateTime := unixToTime.Add(time.Hour * 7)
	return dateTime
}
