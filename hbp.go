package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	blog_config BlogConfig
)

func init() {
	if config_file_exist() == false {
		create_config_file()
		fmt.Fprintln(os.Stderr, "設定ファイルが見つからなかったため,~/.hbpを生成しました.")
		os.Exit(2)
	}
	var err error
	blog_config, err = load_config()
	if err != nil {
		fmt.Fprintln(os.Stderr, "設定ファイルの読み込みに失敗しました.\nホームディレクトリ以下の.hbpファイルを確認してください.")
		os.Exit(1)
	}
}

func main() {
	var (
		fileSpecify     string
		publishSpecify  bool
		categorySpecify string
		versionSpecify  bool
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
	flag.BoolVar(&versionSpecify,
		"v",
		false,
		"バージョンを表示します.")
	flag.Parse()

	if versionSpecify == true {
		fmt.Println("hbp: version 1.0")
		os.Exit(0)
	}

	p := PostConfig{
		SourceFilePath: fileSpecify,
		Publish:        publishSpecify,
		Category:       categorySpecify,
		Author:         blog_config.Hatena_id,
	}

	status, err_msg := post(p)

	if status != 0 {
		fmt.Fprintln(os.Stderr, err_msg)
	}

	os.Exit(status)
}

func post(pcon PostConfig) (status int, err_msg string) {
	var fp *os.File
	fp = os.Stdin
	if pcon.SourceFilePath != "" {
		var err error
		fp, err = os.Open(pcon.SourceFilePath)
		if err != nil {
			return 1, "ファイルオープンエラー."
		}
		defer fp.Close()
	}

	post_xml, xml_create_err := create_post_xml(
		fp,
		blog_config.Hatena_id,
		pcon.Category,
		pcon.Publish)
	if xml_create_err != nil {
		return 1, "投稿内容に問題があります.タイトルの指定などを確認してください."
	}

	api_call_err := call_atom_api(post_xml, blog_config)
	if api_call_err != nil {
		return 1, "APIコールエラー.\n1. 通信状況を確認してください\n2. アカウント設定は正しいですか?  ~/.hbpを確認してください."
	}

	return 0, ""
}
