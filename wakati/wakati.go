package wakati

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/YuheiNakasaka/sayhuuzoku/db"
	"github.com/YuheiNakasaka/sayhuuzoku/scraping"
	"github.com/ikawaha/kagome/tokenizer"
)

// MyToken : token struct
type MyToken struct {
	text string
	pos  int
}

// WakatiShopFileName : wakati shop name file
var WakatiShopFileName = "shoplist_wakati.txt"

// Start : create wakati file
func Start() error {
	fmt.Println("Start creating wakati file")

	// wakti ファイルを開く
	mydb := &db.MyDB{}
	err := mydb.New()
	if err != nil {
		fmt.Println("Failed to open db")
	}
	defer mydb.Connection.Close()

	// ファイルを1行ずつ読み込む準備
	file, err := os.Open(scraping.ShopNameFile)
	if err != nil {
		fmt.Println("Failed to open file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// 読みこんだテキストを貯めるチャンネルと
	// token処理した単語を送るチャンネル
	wg := &sync.WaitGroup{}
	lines := make(chan string)
	values := make(chan MyToken)

	// テキストを処理
	t := tokenizer.New()
	for j := 0; j < 5; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for l := range lines {
				cnt := 1 // 除いたtoken分詰めたposition
				tokens := t.Tokenize(l)
				for _, token := range tokens {
					if token.Class == tokenizer.DUMMY {
						continue
					}
					s := strings.TrimSpace(token.Surface)
					if len(s) > 0 {
						mytoken := MyToken{}
						mytoken.text = s
						mytoken.pos = cnt
						values <- mytoken
						cnt++
					}
				}
			}
		}()
	}

	// fileをガッと読む
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			// 最後の行
			if text == "\n" || text == "" || text == " " {
				return
			}
			lines <- text
		}
		close(lines)
	}()

	go func() {
		wg.Wait()
		close(values)
	}()

	// 100レコードずつ処理(多くしすぎるとtoo many sql variablesのエラーが出る)
	mu := &sync.Mutex{}
	valueQueue := make([]MyToken, 0, 0)
	for v := range values {
		valueQueue = append(valueQueue, v)
		if len(valueQueue) == 100 {
			writeMutex(valueQueue, mydb, mu)
			valueQueue = make([]MyToken, 0, 0)
		}
	}
	writeMutex(valueQueue, mydb, mu)

	fmt.Println("Finish")
	return err
}

// bulk insertする
func writeMutex(values []MyToken, mydb *db.MyDB, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	// valuesをbulk insert用のクエリにする
	valStr := make([]string, 0, 0)
	valArgs := make([]interface{}, 0, 0)
	for _, val := range values {
		if val.text == " " {
			continue
		}
		valStr = append(valStr, "(?, ?)")
		valArgs = append(valArgs, val.text)
		valArgs = append(valArgs, val.pos)
	}
	query := fmt.Sprintf("INSERT INTO wakati_shopname(word, position) values %s", strings.Join(valStr, ","))

	stmt, err := mydb.Connection.Prepare(query)
	if err != nil {
		fmt.Println("Error occured in stmt")
		panic(err)
	}

	_, err = stmt.Exec(valArgs...)
	if err != nil {
		fmt.Println("Error occured in exec")
		panic(err)
	}
}
