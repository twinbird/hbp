# hbp
Hatena Blog Post

はてなブログにコマンドラインから投稿するコマンドラインツールです.
デフォルトでは下書きとして投稿します.

## Usage

標準入力の文字列をはてなブログに投稿します.
最初の行をタイトル,1行空行(改行のみの行)をあけて3行目以降を本文とみなします.

投稿先には設定ファイル~/.hbpの内容を利用します.
~/.hbpが見つからなければテンプレートを生成します.

```sh
# First, you must create configuration file.
$ cat <<EOS
hatena id:Your hatena id
blog id:Your hatena blog id
api key:Your hatena blog api key
EOS > ~/.hbp

$ hbp <<EOS
タイトル

本文
はてなブログにテスト投稿する.
EOS
```

## Options

| Option | Description |
|--------|-------------|
| -f     | (file)指定したファイルの内容で下書きに投稿します. |
| -p     | (publish)下書きではなく公開状態で投稿します. |
| -c     | (category)カテゴリーを指定して投稿します.カテゴリはカンマ(,)区切りで複数指定できます. |
| -v     | (version)バージョン情報を表示します. |

## Example

blog.txtの内容を下書きとして投稿する.

```sh
hbp -f blog.txt
```

blog.txtの内容を公開状態で投稿する.

```sh
hbp -f blog.txt -p
```

blog.txtの内容を下書きとしてカテゴリーfooで投稿する.

```sh
hbp -f blog.txt -c "foo"
```

## Status

| Value | Description |
|-------|-------------|
| 0     | 正常終了    |
| 1     | 異常終了    |
| 2     | 設定ファイルがなかったためテンプレートを生成して終了    |




## Lisence

MIT Lisence
