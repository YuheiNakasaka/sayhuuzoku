package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/YuheiNakasaka/sayhuuzoku/db"
)

// Start : generate shop name
func Start(total int) (string, error) {
	if total < 1 {
		panic("Total count is more than 1.")
	}

	mydb := &db.MyDB{}
	err := mydb.New()
	if err != nil {
		fmt.Println("Failed to connect db")
	}
	defer mydb.Connection.Close()

	words := make([]string, 0, 0)
	for i := 1; i <= total; i++ {
		query := fmt.Sprintf("select * from wakati_shopname where position = %d order by random() limit 1;", i)
		rows, qerr := mydb.Connection.Query(query)
		if qerr != nil {
			fmt.Println("Failed to exexute query.")
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var word string
			var position int
			var createdAt time.Time
			serr := rows.Scan(&id, &word, &position, &createdAt)
			if serr != nil {
				fmt.Println("Failed to fetch row")
			}
			words = append(words, word)
		}
		defer rows.Close()
	}
	return strings.Join(words, ""), err
}
