package main

import (
	"encoding/json"
	"github.com/bot-api/telegram"
	"io/ioutil"
	"strconv"
)

type Question struct {
	ID          int                           `json:"-"`
	Name        string                        `json:"-"`
	Title       string                        `json:"title"`
	Answer      string                        `json:"answer"`
	Items       []*Question                   `json:"items"`
	Parent      *Question                     `json:"-"`
	ReplyMarkup telegram.InlineKeyboardMarkup `json:"-"`
}

type Questions []*Question

func readQuestions(filename string) Questions {
	buf, e := ioutil.ReadFile(filename)
	checkErr(e)

	var root Question
	e = json.Unmarshal(buf, &root)
	checkErr(e)

	items := make(Questions, 0, 32)
	items.append(nil, &root)
	items.update()

	return items
}

func (qs *Questions) append(parent *Question, item *Question) {
	item.ID = len(*qs)
	item.Name = strconv.Itoa(item.ID)
	item.Parent = parent

	*qs = append(*qs, item)
	for _, next := range item.Items {
		qs.append(item, next)
	}
}

func (qs Questions) update() {
	for _, item := range qs {
		itemKeyboard := item
		length := len(itemKeyboard.Items)
		if length == 0 {
			itemKeyboard = item.Parent
			length = len(itemKeyboard.Items)
		}
		length++

		text := make([]string, 0, length)
		data := make([]string, 0, length)
		for _, next := range itemKeyboard.Items {
			text = append(text, next.Title)
			data = append(data, next.Name)
		}

		if itemKeyboard.Parent != nil {
			text = append(text, "\xE2\x86\xA9")
			data = append(data, itemKeyboard.Parent.Name)
		}

		item.ReplyMarkup = telegram.InlineKeyboardMarkup{
			InlineKeyboard: telegram.NewVInlineKeyboard(item.Name+":", text, data),
		}
	}
}
