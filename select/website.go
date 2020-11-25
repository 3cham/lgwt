package website_racer

import (
	"fmt"
	"net/http"
	"time"
)

var (
	tenSecondsTimeout = 10 * time.Second
)

func ConfigurableRacer(fUrl, sUrl string, timeout time.Duration) (winner string, err error) {
	select {
	case <-ping(fUrl):
		return fUrl, nil
	case <-ping(sUrl):
		return sUrl, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout waiting for %s and %s", fUrl, sUrl)
	}
}

func Racer(fUrl, sUrl string) (winner string, err error) {
	return ConfigurableRacer(fUrl, sUrl, tenSecondsTimeout)
}

func ping(url string) chan struct{} {
	channel := make(chan struct{})

	go func() {
		http.Get(url)
		close(channel)
	}()
	return channel
}
