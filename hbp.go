package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
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
		fmt.Println("XML Marshal error:", err)
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
	res, _ := client.Do(req)
	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	if config_file_exist() == false {
		create_config_file()
		fmt.Fprintln(os.Stderr, "設定ファイルが見つからなかったため,~/.hbpを生成しました.")
		os.Exit(2)
	}
	user_configuration = load_config()
}

func main() {
	var (
		fileSpecify     string
		publishSpecify  bool
		categorySpecify string
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`%s
はてなブログ用コマンドラインクライアント.

標準入力をはてなブログへ下書き投稿します.
1行目をタイトル, 3行目以降を本文として扱います.

アカウントの設定は~/.hbpファイルにて行ってください.

[オプション]
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&fileSpecify,
		"f",
		"",
		"指定ファイルを入力として利用します.")
	flag.BoolVar(&publishSpecify,
		"p",
		false,
		"下書きではなく直ちに公開します.")
	flag.StringVar(&categorySpecify,
		"c",
		"",
		"投稿時に付加するカテゴリを指定します.")
	flag.Parse()

	status, err_msg := post(fileSpecify,
		categorySpecify,
		publishSpecify)

	if status != 0 {
		fmt.Fprintln(os.Stderr, err_msg)
	}

	os.Exit(status)
}

func post(fileSpecify, categorySpecify string, publishSpecify bool) (status int, err_msg string) {
	var fp *os.File
	fp = os.Stdin
	if fileSpecify != "" {
		var err error
		fp, err = os.Open(fileSpecify)
		if err != nil {
			return 1, "ファイルオープンエラー."
		}
		defer fp.Close()
	}

	post_xml, xml_create_err := create_post_xml(fp, fileSpecify, categorySpecify, publishSpecify)
	if xml_create_err != nil {
		return 1, "投稿内容に問題があります.タイトルの指定などを確認してください."
	}

	api_call_err := call_atom_api(post_xml)
	if api_call_err != nil {
		return 1, "APIコールエラー.通信状況等を確認してください."
	}

	return 0, ""
}

type UserConfiguration map[string]string

var user_configuration UserConfiguration

func load_config() UserConfiguration {
	user, _ := user.Current()
	config_file_path := (user.HomeDir + "/.hbp")
	fp, err := os.Open(config_file_path)
	if err != nil {
		fmt.Println("設定ファイルが見つかりません.\nホームディレクトリ以下の.hbpファイルを確認してください.")
		os.Exit(1)
	}
	return load_config_file(fp)
}

func load_config_file(fp *os.File) UserConfiguration {
	m := make(UserConfiguration)
	reader := csv.NewReader(fp)
	reader.Comma = ':'
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("設定ファイルの読み込みに失敗しました.\nホームディレクトリ以下の.hbpファイルを確認してください.")
			os.Exit(1)
		}
		m[record[0]] = record[1]
	}
	return m
}

func config_file_exist() bool {
	user, _ := user.Current()
	config_file_path := (user.HomeDir + "/.hbp")
	_, err := os.Stat(config_file_path)
	return err == nil
}

func create_config_file() {
	user, _ := user.Current()
	config_file_path := (user.HomeDir + "/.hbp")

	content := []byte(
		"hatena_id:Your hatena id\n" +
			"blog_id:Your hatena blog id\n" +
			"api_key:Your hatena blog atom api key")
	ioutil.WriteFile(config_file_path, content, os.ModePerm)
}
