package scraping

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ShopListURL : target url
var ShopListURL = "http://fujoho.jp/index.php?p=shop_list&b="

// ShopNameFile : shop name file name
var ShopNameFile = "/scraping/shoplist.txt"

// ShopDicFile : shop dictionary file name
var ShopDicFile = "/scraping/shopdic.txt"

// Start : fetch page and get names
func Start(maxPage int) error {
	fmt.Println("Start scraping")

	absDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return fmt.Errorf("Failed to open file: %v", err)
	}

	// 店名をリストにするファイル
	file, err := os.Create(filepath.Join(absDir, filepath.FromSlash(ShopNameFile)))
	if err != nil {
		return fmt.Errorf("Failed to create file: %v", err)
	}
	defer file.Close()

	// 並列取得する
	maxConnection := make(chan bool, 5)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := 0; i < maxPage; i++ {
		wg.Add(1)
		maxConnection <- true

		go func(page int, mu *sync.Mutex) {
			defer wg.Done()
			time.Sleep(2 * time.Second) // 2秒待つ

			// goqueryでHTML取得
			url := ShopListURL + strconv.Itoa(page)
			doc, scrapingErr := goquery.NewDocument(url)
			if scrapingErr != nil {
				err = fmt.Errorf("Failed to scrape: %v", scrapingErr)
				return
			}
			fmt.Println("Scraping: " + url)

			// 店舗の載ってる範囲をチェック
			if checkValidSite(doc) {
				// 店名を抜き出してファイルに書き出す
				fetchShopName(doc, file, mu)
			}

			<-maxConnection
		}(i, mu)
	}
	wg.Wait()
	fmt.Println("Finish")
	return err
}

func checkValidSite(doc *goquery.Document) bool {
	var t string
	doc.Find(".data-nothing").Each(func(_ int, s *goquery.Selection) {
		t = s.Text()
		fmt.Println(t)
	})
	if t == "" {
		return true
	}
	return false
}

func fetchShopName(doc *goquery.Document, file *os.File, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	doc.Find(".shop-name").Each(func(i int, s *goquery.Selection) {
		shopName := s.Text()
		// 新店舗のmarkは除く
		shopName = strings.Replace(shopName, "New!", "", 1)
		// カッコの補足は除く
		rep1 := regexp.MustCompile(`[\(|（].+[\)|）]`)
		shopName = rep1.ReplaceAllString(shopName, "")

		// ファイルに書き込む
		file.Write(([]byte)(shopName + "\n"))
		// fmt.Printf("Result %d: %s\n", i, shopName)
	})
}
