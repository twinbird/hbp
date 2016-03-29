package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	HATENA_BASE_URL = "https://blog.hatena.ne.jp/"
)

type Entries struct {
	XMLName  xml.Name `xml:"entry"`
	XMLns    string   `xml:"xmlns,attr"`
	XMLnsApp string   `xml:"xmlns:app,attr"`
	Title    string   `xml:"title"`
	Author   Author   `xml:"author"`
	Content  Content  `xml:"content"`
	Updated  string   `xml:"updated"`
	Category Category `xml:"category"`
	App      App
}

type Author struct {
	Name string `xml:"name"`
}

type Content struct {
	XMLName      xml.Name `xml:"content"`
	ContentValue string   `xml:",chardata"`
	Type         string   `xml:"type,attr"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	Term    string   `xml:"term,attr"`
}

type App struct {
	XMLName xml.Name `xml:"app:control"`
	Draft   string   `xml:"app:draft"`
}

func create_post_xml(fp *os.File, fileSpecify, categorySpecify string, publishSpecify bool) (post_xml string, ret_err error) {
	reader := bufio.NewReader(fp)

	title, _ := reader.ReadString('\n')
	title = strings.TrimRight(title, "\n")

	reader.ReadString('\n')

	content := ""
	buf := ""
	var err error
	for {
		buf, err = reader.ReadString('\n')
		content += buf
		if err == io.EOF {
			break
		}
	}

	draft := true
	if publishSpecify == true {
		draft = false
	}

	build_xml := build_post_xml(string(title), user_configuration["hatena_id"], content, string(categorySpecify), draft)

	return build_xml, nil
}

func build_post_xml(title, author, contents, category string, draft bool) string {
	draft_str := "yes"
	if draft == false {
		draft_str = "no"
	}
	v := &Entries{
		XMLns:    "http://www.w3.org/2005/Atom",
		XMLnsApp: "http://www.w3.org/2007/app",
		Title:    title,
		Author:   Author{Name: author},
		Content:  Content{ContentValue: contents, Type: "text/plain"},
		Category: Category{Term: category},
		App:      App{Draft: draft_str},
	}
	output, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "投稿用XMLの生成に失敗しました:", err)
		return "error"
	}

	return xml.Header + string(output)
}

func draft_post_url() string {
	url := fmt.Sprintf("%s%s/%s/atom/entry",
		HATENA_BASE_URL,
		user_configuration["hatena_id"],
		user_configuration["blog_id"])
	return url
}

func call_atom_api(xml string) error {
	draft_post_url := draft_post_url()
	req, _ := http.NewRequest(
		"POST",
		draft_post_url,
		bytes.NewBuffer([]byte(xml)))
	req.SetBasicAuth(user_configuration["hatena_id"], user_configuration["api_key"])
	req.Header.Set("Content-Type", "application/atomsvc+xml; charset=utf-8")

	client := new(http.Client)
	res, req_err := client.Do(req)
	defer res.Body.Close()
	if req_err != nil || res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "HTTPリクエストエラー.\nアカウント設定は正しいですか?\n~/.hbpを確認してください.\n")
		os.Exit(1)
	}

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}
