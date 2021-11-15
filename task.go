package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-sdk/lib/codec/json"
	"github.com/go-sdk/lib/cron"
	"github.com/go-sdk/lib/log"
)

type Task struct {
	Name       string            `json:"name"`
	Url        string            `json:"url"`
	Body       string            `json:"body"`
	Cron       string            `json:"cron"`
	Timezone   string            `json:"timezone"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Timeout    int64             `json:"timeout"`
	NoBody     bool              `json:"no_body"`
	Decryption string            `json:"decryption"`
	HTTPProxy  string            `json:"http_proxy"`
	HTTPSProxy string            `json:"https_proxy"`
}

func StartTask() error {
	if err := initConfig(); err != nil {
		return err
	}

	c := cron.Default(nil)
	for i := range config.Tasks {
		i := i
		t := config.Tasks[i]
		c.Add(t.Cron, t.Name, func() { do(i, log.WithFields(log.Fields{"name": t.Name})) })
	}
	c.Start()

	return nil
}

func initConfig() error {
	bs, err := ioutil.ReadFile(config.Path)
	if err != nil {
		return fmt.Errorf("read file fail, %v", err)
	}

	err = json.Unmarshal(bs, &config)
	if err != nil {
		return fmt.Errorf("decode file fail, %v", err)
	}

	if len(config.Tasks) == 0 {
		return fmt.Errorf("task empty")
	}

	for i, task := range config.Tasks {
		if task.Cron == "" {
			config.Tasks[i].Cron = "* * * * * *"
		}
		if task.Timezone == "" {
			config.Tasks[i].Timezone = time.UTC.String()
		}
		if task.Method == "" {
			config.Tasks[i].Method = "GET"
		}
		if task.Timeout <= 0 || task.Timeout > 60*30 {
			config.Tasks[i].Timeout = 60
		}
	}

	log.Debugf("config: %s", json.MustMarshal(config))

	return nil
}

func do(i int, l *log.Entry) {
	task := config.Tasks[i]

	req, err := http.NewRequest(task.Method, task.Url, bytes.NewReader([]byte(task.Body)))
	if err != nil {
		l.WithField("err", err).Errorf("build request fail")
		return
	}

	for k, v := range task.Headers {
		req.Header.Add(k, v)
	}

	now := time.Now()

	resp, err := (&http.Client{Transport: transport(task.HTTPProxy, task.HTTPSProxy), Timeout: time.Duration(task.Timeout) * time.Second}).Do(req)
	if err != nil {
		l.WithField("err", err).Errorf("send request fail")
		return
	}
	defer resp.Body.Close()

	took := time.Since(now).Milliseconds()
	l.Debugf("took %dms\n%s", took, dump(task, req, resp))
	l.Infof("end, took %dms, status %s", took, resp.Status)
}

func transport(hp, hsp string) http.RoundTripper {
	var hpu, hspu *url.URL
	if hp != "" {
		hpu, _ = url.Parse(hp)
	}
	if hsp != "" {
		hspu, _ = url.Parse(hsp)
	}
	proxy := func(req *http.Request) (*url.URL, error) {
		var u *url.URL
		if req.URL.Scheme == "https" {
			u = hspu
		} else {
			u = hpu
		}
		return u, nil
	}
	return &http.Transport{Proxy: proxy}
}

func dump(task Task, req *http.Request, resp *http.Response) string {
	bs, _ := ioutil.ReadAll(resp.Body)

	sb := strings.Builder{}

	// Request
	sb.WriteString(req.Method + " " + req.URL.String() + " " + req.Proto + "\n")
	sb.WriteString("Host: " + req.Host + "\n")
	for k, v := range req.Header {
		sb.WriteString(k + ": " + strings.Join(v, " ") + "\n")
	}
	if task.Body == "" {
		task.Body = "<empty>"
	}
	sb.WriteString("\n" + task.Body + "\n")

	sb.WriteString("--------------------------------------------------------------------------------\n")

	// Response
	sb.WriteString(resp.Proto + " " + resp.Status + "\n")
	for k, v := range resp.Header {
		sb.WriteString(k + ": " + strings.Join(v, " ") + "\n")
	}
	if len(bs) > 0 && !task.NoBody {
		var (
			str string
			err error
		)
		sb.WriteString("\n")
		switch strings.ToLower(task.Decryption) {
		case "unicode":
			str, err = strconv.Unquote(strings.Replace(strconv.Quote(string(bs)), `\\u`, `\u`, -1))
		default:
		}
		if err == nil {
			sb.WriteString(str)
		} else {
			sb.WriteString(string(bs))
		}
	}

	return sb.String()
}
