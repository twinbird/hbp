package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

type PostConfig struct {
	SourceFilePath string
	Publish        bool
	Category       string
	Author         string
}

func create_post_xml(
	fp *os.File,
	author, categorySpecify string,
	publishSpecify bool) (post_xml string, ret_err error) {

	reader := bufio.NewReader(fp)

	// First Row is Title
	title, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	title = strings.TrimRight(title, "\n")

	// Second Row is Skip
	reader.ReadString('\n')

	// Third Row or Later is Contents
	content := ""
	buf := ""
	for {
		buf, ret_err = reader.ReadString('\n')
		content += buf
		if ret_err == io.EOF {
			break
		} else if ret_err != nil {
			return "", ret_err
		}
	}

	return build_post_xml(string(title),
		author,
		content,
		string(categorySpecify),
		publishSpecify)
}

func build_post_xml(title, author, contents, category string, publish bool) (build_xml string, err error) {
	draft_str := "yes"
	if publish == true {
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
	output, msl_err := xml.MarshalIndent(v, "", "    ")
	if msl_err != nil {
		return "", msl_err
	}

	return xml.Header + string(output), nil
}

func call_atom_api(xml string, config BlogConfig) error {
	draft_post_url := fmt.Sprintf("%s%s/%s/atom/entry",
		config.Base_url,
		config.Hatena_id,
		config.Blog_id)

	req, _ := http.NewRequest(
		"POST",
		draft_post_url,
		bytes.NewBuffer([]byte(xml)))
	req.SetBasicAuth(config.Hatena_id, config.Api_key)
	req.Header.Set("Content-Type", "application/atomsvc+xml; charset=utf-8")

	client := new(http.Client)
	res, req_err := client.Do(req)
	defer res.Body.Close()

	if req_err != nil {
		return req_err
	}

	if res.StatusCode != http.StatusCreated {
		errors.New("HTTP But Response")
	}

	return nil
}
