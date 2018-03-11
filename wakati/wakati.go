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

	// ファイルを1行ずつ読み込む
	file, err := os.Open(scraping.ShopNameFile)
	if err != nil {
		fmt.Println("Failed to open file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// 店名を分かち書きしてファイルに書き出し
	maxConnection := make(chan bool, 20)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	values := make([]MyToken, 0, 0)
	t := tokenizer.New()
	for scanner.Scan() {
		wg.Add(1)
		maxConnection <- true

		go func() {
			defer wg.Done()
			text := scanner.Text()
			// 最後の行
			if text == "\n" || text == "" || text == " " {
				<-maxConnection
				return
			}
			fmt.Println(text)

			cnt := 1 // 除いたtoken分詰めたposition
			tokens := t.Tokenize(text)
			for _, token := range tokens {
				if token.Class == tokenizer.DUMMY {
					continue
				}
				s := strings.TrimSpace(token.Surface)
				if len(s) > 0 {
					mytoken := MyToken{}
					mytoken.text = s
					mytoken.pos = cnt
					values = append(values, mytoken)
					cnt++
				}
			}
			<-maxConnection
		}()
	}
	wg.Wait()

	// 100レコードずつ処理(多くしすぎるとtoo many sql variablesのエラーが出る)
	idx := 0
	for i := range values {
		if i%100 == 0 {
			writeMutex(values[i:i+100], mydb, mu)
			idx = i
		}
	}
	writeMutex(values[idx:], mydb, mu)

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
