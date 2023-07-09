package util

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var r = resty.New()

func RGet(url string) string {
	v, err := r.R().Get(url)
	if err != nil {
		log.Errorf("rest get error: %v", err)
		return ""
	}
	return v.String()
}

func RPost(url string, body string) string {
	v, err := r.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"testuser", "password":"testpass"}`).
		Post(url)
	if err != nil {
		log.Errorf("rest post error: %v", err)
		return ""
	}
	return v.String()
}

func RPut(url string, body string) string {
	v, err := r.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"testuser", "password":"testpass"}`).
		Put(url)
	if err != nil {
		log.Errorf("rest put error: %v", err)
		return ""
	}
	return v.String()
}

func RPatch(url string, body string) string {
	v, err := r.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"testuser", "password":"testpass"}`).
		Patch(url)
	if err != nil {
		log.Errorf("rest patch error: %v", err)
		return ""
	}
	return v.String()
}

func RDelete(url string) {
	_, err := r.R().Delete(url)
	if err != nil {
		log.Errorf("rest delete error: %v", err)
	}
}
