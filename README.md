# 9cc-go

## C言語からアセンブリを作成するコンパイラ
### 進捗状況
- 四則演算と比較関数(<, >, ==, !=, <= ,>=)と括弧( ) の実装
- 複数文字変数の参照と代入
- for, if, while 文の実装
- block の実装 { }
- 6引数までの関数実装、呼び出し、RSP 16の倍数対応(ABI)

## Usage
``` consolev
Usage:
  9cc-go [flags]

Flags:
  -h, --help            help for 9cc-go
  -o, --output string   output file name (default "output.s")
      --show-tokenize   show tokenize
  -t, --toggle          Help message for toggle
      --tokenize        only tokenize

```

## Reference
- 低レイヤを知りたい人のためのCコンパイラ作成入門 : https://www.sigbus.info/compilerbook
- x86アセンブリ言語での関数コール : https://vanya.jp.net/os/x86call/
