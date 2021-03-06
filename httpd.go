package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"log"
	"math"
	"time"
)

type Httpd struct {
	Host string
	Port int
}

type Param struct {
	Channel    string   `form:"channel" binding:"required"`
	Message    string   `form:"message"`
	Name       string   `form:"name"`
	Icon       string   `form:"icon"`
	Fallback   string   `form:"fallback"`
	Color      string   `form:"color"`
	Pretext    string   `form:"pretext"`
	AuthorName string   `form:"author_name"`
	AuthorLink string   `form:"author_link"`
	AuthorIcon string   `form:"author_icon"`
	Title      string   `form:"title"`
	TitleLink  string   `form:"title_link"`
	Text       string   `form:"text"`
	FieldTitle []string `form:"field_title[]"`
	FieldValue []string `form:"field_value[]"`
	FieldShort []bool   `form:"field_short[]"`
	ImageURL   string   `form:"image_url"`
	Manual     bool     `form:"manual"`
	PostAt     int64    `form:"post_at"`
}

func NewHttpd(host string, port int) *Httpd {
	return &Httpd{
		Host: host,
		Port: port,
	}
}

func (h *Httpd) Run() {
	m := martini.Classic()
	m.Get("/", func() string { return "Hello, I'm Takosan!!1" })
	m.Post("/notice", binding.Bind(Param{}), messageHandler)
	m.Post("/privmsg", binding.Bind(Param{}), messageHandler)
	m.RunOnAddr(fmt.Sprintf("%s:%d", h.Host, h.Port))
}

func messageHandler(p Param) (int, string) {
	ch := make(chan error, 1)

	if p.PostAt > 0 {
		diff := p.PostAt - time.Now().Unix()
		delay := int64(math.Max(float64(diff), 0))

		return sendLater(p, delay, ch)
	} else {
		return sendNow(p, ch)
	}
}

func sendNow(p Param, ch chan error) (int, string) {
	go MessageBus.Publish(NewMessage(p, ch), 0)
	err := <-ch

	if err != nil {
		message := fmt.Sprintf("Failed to send message to %s: %s\n", p.Channel, err)
		log.Printf(fmt.Sprintf("[error] %s", message))
		return 400, message
	} else {
		return 200, fmt.Sprintf("Message sent successfully to %s", p.Channel)
	}
}

func sendLater(p Param, delay int64, ch chan error) (int, string) {
	go MessageBus.Publish(NewMessage(p, ch), delay)

	return 200, fmt.Sprintf("Message accepted and will be sent after %d seconds", delay)
}
