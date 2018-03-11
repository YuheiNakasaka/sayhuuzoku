# say-huuzoku

## Description
The sayhuuzoku is a library to generate a shop name like 風俗店(huuzoku-shop).

## Example
```
$ go run sayhuuzoku.go generate -c 2
月淫乱

$ go run sayhuuzoku.go generate -c 3
セレブサークル月ちゃり

$ go run sayhuuzoku.go generate -c 4
INO-遊園PLAYSTAGE
```

## Installation
```
go get -u github.com/YuheiNakasaka/sayhuuzoku
```

## Usage

```
$ sayhuuzoku -h
NAME:
   sayhuuzoku - A new cli application

USAGE:
   sayhuuzoku [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     init, i      Init database
     scraping, s  Fetch shop name from http://fujoho.jp/index.php?p=shop_list
     wakati, w    Create wakati data from shoplist file
     generate, g  Generate shop name like huuzoku
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## License
The library is available as open source under the terms of the MIT License.
