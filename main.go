package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"
	"github.com/robfig/cron/v3"
)

type Config struct {
	Tasks []Task `json:"tasks"`

	Path string `json:"-"`
	Help bool   `json:"-"`
}

type Task struct {
	Name       string            `json:"name"`
	Url        string            `json:"url"`
	Body       string            `json:"body"`
	Cron       string            `json:"cron"`
	Timezone   string            `json:"timezone"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Timeout    int64             `json:"timeout"`
	HTTPProxy  string            `json:"http_proxy"`
	HTTPSProxy string            `json:"https_proxy"`
}

var (
	config = &Config{}
)

func init() {
	flag.StringVar(&config.Path, "config", "config.json", "config")
	flag.BoolVar(&config.Help, "help", false, "instructions for use")
	flag.Parse()

	if config.Help {
		flag.Usage()
		os.Exit(0)
	}

	bs, err := ioutil.ReadFile(config.Path)
	if err != nil {
		logx.WithField("err", err).Fatal("[init] read file fail")
	}

	err = json.Unmarshal(bs, &config)
	if err != nil {
		logx.WithField("err", err).Fatal("[init] decode file fail")
	}

	if len(config.Tasks) == 0 {
		logx.Fatal("[init] task empty")
	}

	for i, task := range config.Tasks {
		if task.Cron == "" {
			config.Tasks[i].Cron = "* * * * *"
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

	logx.Debugf("[init] config: %s", json.MustMarshal(config))
}

func main() {
	for i := range config.Tasks {
		go handle(i)
	}
	select {}
}

func handle(i int) {
	task := config.Tasks[i]

	l, err := time.LoadLocation(task.Timezone)
	if err != nil {
		logx.WithField("err", err).Fatalf("[%s:%d] timezone parse fail", task.Name, i)
	}
	s, err := cron.ParseStandard(task.Cron)
	if err != nil {
		logx.WithField("err", err).Fatalf("[%s:%d] cron parse fail", task.Name, i)
	}

	c := cron.New(cron.WithLocation(l), cron.WithLogger(cron.VerbosePrintfLogger(&Log{i: i, name: task.Name})))
	c.Schedule(s, cron.FuncJob(func() { do(i) }))
	c.Start()
}

func do(i int) {
	task := config.Tasks[i]

	logx.Infof("[%s:%d] begin", task.Name, i)

	req, err := http.NewRequest(task.Method, task.Url, bytes.NewReader([]byte(task.Body)))
	if err != nil {
		logx.WithField("err", err).Errorf("[%s:%d] build request fail", task.Name, i)
		return
	}

	for k, v := range task.Headers {
		req.Header.Add(k, v)
	}

	now := time.Now()

	resp, err := (&http.Client{Transport: transport(task.HTTPProxy, task.HTTPSProxy), Timeout: time.Duration(task.Timeout) * time.Second}).Do(req)
	if err != nil {
		logx.WithField("err", err).Errorf("[%s:%d] send request fail", task.Name, i)
		return
	}
	defer resp.Body.Close()

	took := time.Since(now).Milliseconds()
	logx.Debugf("[%s:%d] took %dms\n%s", task.Name, i, took, dump(task, req, resp))
	logx.Infof("[%s:%d] end, took %dms, status %s", task.Name, i, took, resp.Status)
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
	if len(bs) > 0 {
		sb.WriteString("\n" + hex.EncodeToString(bs))
	}

	return sb.String()
}
