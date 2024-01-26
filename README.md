# msys2wtp

## Windows TerminalでMSYS2関係の環境を実行するための記述を生成するツール

Windows TerminalでMSYS2関係の環境を実行するには、設定画面から新しいプロファイルを一つずつ追加していく方法と、`settings.json` をエディタで開いて `profiles` にまとめて追加する方法の2種類があります。これは後者の方法で行うための記述内容(json形式の一部分)を自動生成するツールです。

## 使用方法

``` shell
$ ./msys2wtp.exe

Options:
  -c    output to clipboard
  -d string
        starting directory (default "C:\\msys64\\home\\%USERNAME%")
  -gfw
        adding a entry of Git for Windows
  -i string
        MSYS2 install path (default "C:/msys64")
  -s string
        shell (bash/zsh) (default "bash")
  -t string
        types of msys2 (default "msys2,ucrt64,mingw32,mingw64,clang64,clang32,clangarm64")
```

- `-c` : `-c` オプションがあると、クリップボードに出力します。`-c` オプションがない場合は、標準出力に出力します。
- `-d` : 起動時のディレクトリを指定します
- `-i` : MSYS2をインストールしたパスを指定します
- `-s` : 起動するシェルを指定します
- `-t` : MSYS2の環境をカンマ区切りで指定します
- `-gfw` : Git for Windows用のプロファイルも生成します

`settings.json` にカットアンドペースしやすいように、 `profiles` の `list` 配列要素の部分だけを出力します。`settings.json` にペーストする際は、前後のカンマの有無に注意してください。

## 出力サンプル

output.txt は出力サンプルです。

## ライセンス

MIT
