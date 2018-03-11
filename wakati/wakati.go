package wakati

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/YuheiNakasaka/sayhuuzoku/scraping"
	"github.com/ikawaha/kagome/tokenizer"
)

// WakatiShopFileName : wakati shop name file
var WakatiShopFileName = "shoplist_wakati.txt"

// Start : create wakati file
func Start() error {
	fmt.Println("Start creating wakati file")

	// wakti ファイルを開く
	wakatiFile, err := os.Create(WakatiShopFileName)
	if err != nil {
		fmt.Println("Failed to create file")
	}
	defer wakatiFile.Close()

	// ファイルを1行ずつ読み込む
	file, err := os.Open(scraping.ShopNameFile)
	if err != nil {
		fmt.Println("Failed to open file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// 店名を分かち書きしてファイルに書き出し
	maxConnection := make(chan bool, 10)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	t := tokenizer.New()
	for scanner.Scan() {
		wg.Add(1)
		maxConnection <- true

		go func(mu *sync.Mutex) {
			defer wg.Done()
			text := scanner.Text()
			fmt.Printf("Tokenize: %s\n", scanner.Text())
			tokens := t.Tokenize(text)
			for _, token := range tokens {
				if token.Class == tokenizer.DUMMY {
					continue
				}
				for _, s := range strings.Split(token.Surface, ",") {
					s = strings.TrimSpace(s)
					if len(s) > 0 {
						writeMutex(s, wakatiFile, mu)
					}
				}
			}
			<-maxConnection
		}(mu)
	}
	wg.Wait()
	fmt.Println("Finish")
	return err
}

func writeMutex(s string, file *os.File, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	file.Write(([]byte)(s + "\n"))
}
