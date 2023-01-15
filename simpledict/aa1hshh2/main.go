package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}
type BaiduRequest struct {
	q     string
	from  string
	to    string
	appid string
	salt  string
	sign  string
}
type BaiduResponse struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func constructbaidurequest(word string) BaiduRequest {
	request := BaiduRequest{q: word, from: "en", to: "zh", salt: "1435660288"}
	request.appid = "" // appid here
	pass := ""         // password here
	data := []byte(request.appid + request.q + request.salt + pass)
	//fmt.Printf("%x", md5.Sum(data))
	request.sign = fmt.Sprintf("%x", md5.Sum(data))
	return request
}
func queryfrombaiduapi(word string, ch chan string) {
	client := &http.Client{}
	request := constructbaidurequest(word)
	//fmt.Println(request)
	url := "http://api.fanyi.baidu.com/api/trans/vip/translate?q=" + word + "&from=" + request.from + "&to=" + request.to +
		"&appid=" + request.appid + "&salt=" + request.salt + "&sign=" + request.sign
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body) //update ioutil.ReadAll() to io.ReadAll()
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var baiduResponse BaiduResponse
	err = json.Unmarshal(bodyText, &baiduResponse)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(word + " ' meaning:")
	rst := "baidu result:"
	for _, item := range baiduResponse.TransResult {
		rst += (item.Dst + "\n")
	}
	ch <- rst
}
func queryfromcaiyunai(word string, ch chan string) {
	client := &http.Client{}
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	rst := "caiyun result:" + "\n"
	rst += fmt.Sprintf("%s", "UK:"+dictResponse.Dictionary.Prons.En+"US:"+dictResponse.Dictionary.Prons.EnUs)
	rst += "\n"
	for _, item := range dictResponse.Dictionary.Explanations {
		rst += (item + "\n")
	}
	ch <- rst
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	ch1 := make(chan string)
	ch2 := make(chan string)

	go queryfrombaiduapi(word, ch1)
	go queryfromcaiyunai(word, ch2)
	fmt.Println(word + " ' meaning:")
	for i := 0; i < 2; i++ {
		select {
		case baidures := <-ch1:
			fmt.Println(baidures)
		case caiyunres := <-ch2:
			fmt.Println(caiyunres)
		}

	}

}
