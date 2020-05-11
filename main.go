package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"
	"github.com/robfig/cron/v3"
)

type Config struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Body     string            `json:"body"`
	Cron     string            `json:"cron"`
	Timezone string            `json:"timezone"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Timeout  int64             `json:"timeout"`
}

var (
	config = Config{}

	c string
	d bool
)

func init() {
	flag.StringVar(&c, "config", "config.json", "config")
	flag.BoolVar(&d, "debug", strings.ToLower(os.Getenv("DEBUG")) == "true", "debug")
	flag.Parse()

	if !d {
		logx.SetLevel(logx.InfoLevel)
	}

	bs, err := ioutil.ReadFile(c)
	if err != nil {
		exit("[init] read file fail", err)
	}

	err = json.Unmarshal(bs, &config)
	if err != nil {
		exit("[init] decode file fail", err)
	}

	if len(config.Tasks) == 0 {
		exit("[init] task empty", nil)
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
		exit("[%s:%d] timezone parse fail", err, task.Name, i)
	}
	s, err := cron.ParseStandard(task.Cron)
	if err != nil {
		exit("[%s:%d] cron parse fail", err, task.Name, i)
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
		exit("[%s:%d] build request fail", err, task.Name, i)
		return
	}

	for k, v := range task.Headers {
		req.Header.Add(k, v)
	}

	now := time.Now()

	resp, err := (&http.Client{Timeout: time.Duration(task.Timeout) * time.Second}).Do(req)
	if err != nil {
		exit("[%s:%d] send request fail", err, task.Name, i)
		return
	}
	defer resp.Body.Close()

	took := time.Since(now).Milliseconds()
	logx.Debugf("[%s:%d] took %dms\n%s", task.Name, i, took, dump(task, req, resp))
	logx.Infof("[%s:%d] end, took %dms, status %s", task.Name, i, took, resp.Status)
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
