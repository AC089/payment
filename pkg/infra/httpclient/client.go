/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-16 19:32:47
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-16 19:32:50
 */
package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var(
	tr = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client = &http.Client{
		Transport: tr,
		Timeout: 5 * time.Second,
	}
)

func GetInstance() *http.Client {
	return client
}

//发送GET请求
//url:请求地址
//response:请求返回的内容
func Get(url string) (response string, err error) {
    resp, err := client.Get(url)
    if err != nil {
        log.WithFields(log.Fields{
			"err": err,
			"url": url,
		}).Error("http client get error")
		return "", err
    }
	defer resp.Body.Close()
    var buffer [512]byte
    result := bytes.NewBuffer(nil)
    for {
        n, err := resp.Body.Read(buffer[0:])
        result.Write(buffer[0:n])
        if err != nil && err == io.EOF {
            break
        } else if err != nil {
            log.WithFields(log.Fields{
				"err": err,
				"url": url,
			}).Error("http client get error")
			return "", err
        }
    }

    response = result.String()
    return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func Post(url string, data interface{}, contentType string) (content string, err error) {
    jsonStr, _ := json.Marshal(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Add("content-type", contentType)
    if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"url": url,
		}).Error("http client get error")
		return "", err
    }
    defer req.Body.Close()

    resp, err := client.Do(req)
    if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"url": url,
		}).Error("http client get error")
		return "", err
    }
    defer resp.Body.Close()

    result, _ := ioutil.ReadAll(resp.Body)
    content = string(result)
    return
}