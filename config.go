package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
)

const (
	HATENA_BASE_URL = "https://blog.hatena.ne.jp/"
)

type BlogConfig struct {
	Base_url  string
	Blog_id   string
	Hatena_id string
	Api_key   string
}

func load_config() (config BlogConfig, ret_err error) {
	user, _ := user.Current()
	config_file_path := (user.HomeDir + "/.hbp")
	fp, err := os.Open(config_file_path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "設定ファイルが見つかりません.\nホームディレクトリ以下の.hbpファイルを確認してください.")
		os.Exit(1)
	}
	return load_config_file(fp)
}

func load_config_file(fp *os.File) (config BlogConfig, ret_err error) {
	config.Base_url = HATENA_BASE_URL
	reader := csv.NewReader(fp)
	reader.Comma = ':'
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return config, err
		}
		switch record[0] {
		case "hatena_id":
			config.Hatena_id = record[1]
		case "blog_id":
			config.Blog_id = record[1]
		case "api_key":
			config.Api_key = record[1]
		}
	}
	return config, nil
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

	content := []byte(`
hatena_id:Your hatena id
blog_id:Your hatena blog id
api_key:Your hatena blog atom api key`)
	ioutil.WriteFile(config_file_path, content, os.ModePerm)
}
