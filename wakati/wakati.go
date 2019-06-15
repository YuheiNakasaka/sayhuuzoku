package wakati

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/YuheiNakasaka/sayhuuzoku/db"
	"github.com/YuheiNakasaka/sayhuuzoku/scraping"
	"github.com/ikawaha/kagome.ipadic/tokenizer"
)

// MyToken : token struct
type MyToken struct {
	text string
	pos  int
}

// Start : create wakati file
func Start() error {
	start := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// wakti ファイルを開く
	db.InitDB = true // dbファイルを初期化する
	mydb := &db.MyDB{}
	err := mydb.New()
	if err != nil {
		return fmt.Errorf("Failed to open db: %v", err)
	}
	defer mydb.Connection.Close()

	// ファイルを1行ずつ読み込む準備
	absDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return fmt.Errorf("Failed to open file: %v", err)
	}
	file, err := os.Open(filepath.Join(absDir, filepath.FromSlash(scraping.ShopNameFile)))
	if err != nil {
		return fmt.Errorf("Failed to open shop name file: %v", err)
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

	dic, err := os.Open(absDir + scraping.ShopDicFile)
	if err != nil {
		return fmt.Errorf("Failed to open user dictionary file: %v", err)
	}
	defer dic.Close()
	userDicRec, err := tokenizer.NewUserDicRecords(dic)
	if err != nil {
		return fmt.Errorf("Failed to create user dictionary record: %v", err)
	}
	userDic, err := userDicRec.NewUserDic()
	if err != nil {
		return fmt.Errorf("Failed to create user dictionary: %v", err)
	}
	t.SetUserDic(userDic)
	for j := 0; j < 100; j++ {
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
					s, nerr := normalize(token)
					if nerr != nil {
						continue
					}
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

	// 多くしすぎるとtoo many sql variablesのエラーが出るのでぎりぎりまで
	mu := &sync.Mutex{}
	valueQueue := make([]MyToken, 0, 0)
	for v := range values {
		valueQueue = append(valueQueue, v)
		if len(valueQueue) == 200 {
			if err := writeMutex(valueQueue, mydb, mu); err != nil {
				return err
			}
			valueQueue = make([]MyToken, 0, 0)
		}
	}
	writeMutex(valueQueue, mydb, mu)

	end := time.Now()
	fmt.Printf("Finish: %f秒\n", (end.Sub(start)).Seconds())
	return nil
}

// bulk insertする
func writeMutex(values []MyToken, mydb *db.MyDB, mu *sync.Mutex) error {
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
		return fmt.Errorf("Error occured in stmt: %v", err)
	}

	_, err = stmt.Exec(valArgs...)
	if err != nil {
		return fmt.Errorf("Error occured in exec: %v", err)
	}
	return nil
}

// 店名の正規化をする
func normalize(token tokenizer.Token) (string, error) {
	var err error
	s := strings.TrimSpace(token.Surface)
	features := token.Features()

	for _, f := range features {
		if f == "空白" || f == "助詞" || f == "助動詞" || f == "サ変接続" || f == "括弧開" || f == "括弧閉" || f == "句点" || f == "地域" {
			err = fmt.Errorf("Invalid style word: %s", f)
			return "", err
		}
		if s == "-" || s == "~" || s == "～" || s == "ー" || s == "店" ||
			s == "." || s == "！" || s == "・" || s == "っ" || s == "s" || s == "ぽ" ||
			s == "…" || s == "？" || s == "、" || s == "倶楽部" || s == "club" || s == "CLUB" ||
			s == "クラブ" || s == "Club" || s == "＆" || s == "☆" || s == "お" {
			err = fmt.Errorf("Stop word: %s", f)
			return "", err
		}
		if f == "動詞" && len(s) == 1 {
			err = fmt.Errorf("Unusal word: %s %s", f, s)
			return "", err
		}
		if f != "名詞" && len(s) == 1 {
			err = fmt.Errorf("Unusal word: %s %s", f, s)
			return "", err
		}
	}
	//fmt.Println(s, token.Features())
	return s, err
}
