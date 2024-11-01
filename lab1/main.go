package main

import (
	"bytes"
	"encoding/json"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
)

var authorBlacklist = []string{
	"author", "Author", "writer", "CEO",
}

type ArticleRecordIn struct {
	Category string `json:"category"`
	Authors  string `json:"authors"`
}

func (r ArticleRecordIn) Out() ([]ArticleRecordOut, bool) {
	r.Authors = strings.ReplaceAll(r.Authors, ", AP", "")
	r.Authors = strings.ReplaceAll(r.Authors, ", M.D", "")

	if len(r.Authors) == 0 || len(r.Category) == 0 {
		return nil, false
	}

	out := []ArticleRecordOut{}

	for _, author := range strings.Split(r.Authors, " and ") {
		index := strings.Index(author, ", ")
		if index > 0 {
			author = author[:index]
		}

		if len(author) > 30 || lo.Contains(authorBlacklist, author) {
			continue
		}

		out = append(out, ArticleRecordOut{
			Category: r.Category,
			Author:   author,
		})
	}

	return out, true
}

type ArticleRecordOut struct {
	Category string `json:"category"`
	Author   string `json:"author"`
}

func main() {
	buf, err := os.ReadFile("./lab1/News_Category_Dataset_v3.json")
	if err != nil {
		log.Fatal(err)
	}

	buf = append([]byte("["), bytes.Join(bytes.Split(buf, []byte{'}'}), []byte("},"))...)
	buf = buf[:len(buf)-2]
	buf = append(buf, byte(']'))

	in := []ArticleRecordIn{}
	out := []ArticleRecordOut{}

	if err := json.Unmarshal(buf, &in); err != nil {
		log.Fatal(err)
	}

	for _, item := range in {
		if o, ok := item.Out(); ok {
			out = append(out, o...)
		}
	}

	out = lo.Uniq(out)

	buf, err = json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(out))

	log.Print(os.WriteFile("./lab1/filtered.json", buf, 0666))
}
