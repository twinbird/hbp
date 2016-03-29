package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
)

type UserConfiguration map[string]string

var user_configuration UserConfiguration

func load_config() UserConfiguration {
	user, _ := user.Current()
	config_file_path := (user.HomeDir + "/.hbp")
	fp, err := os.Open(config_file_path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "設定ファイルが見つかりません.\nホームディレクトリ以下の.hbpファイルを確認してください.")
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
			fmt.Fprintln(os.Stderr, "設定ファイルの読み込みに失敗しました.\nホームディレクトリ以下の.hbpファイルを確認してください.")
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
