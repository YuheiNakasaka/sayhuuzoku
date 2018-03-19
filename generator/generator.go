package generator

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/YuheiNakasaka/sayhuuzoku/db"
)

// Start : generate shop name
func Start(total int) (string, error) {
	if total < 1 {
		return "", errors.New("Total count is more than 1")
	}

	mydb := &db.MyDB{}
	err := mydb.New()
	if err != nil {
		return "", fmt.Errorf("Failed to connect db: %v", err)
	}
	defer mydb.Connection.Close()

	words := make([]string, 0, 0)
	for i := 1; i <= total; i++ {
		query := fmt.Sprintf("select * from wakati_shopname where length(word) > 1 and position = %d group by word order by random() limit 1;", i)
		rows, qerr := mydb.Connection.Query(query)
		if qerr != nil {
			return "", fmt.Errorf("Failed to execute query: %v", err)
		}

		for rows.Next() {
			var id int
			var word string
			var position int
			var createdAt time.Time
			serr := rows.Scan(&id, &word, &position, &createdAt)
			if serr != nil {
				rows.Close()
				return "", fmt.Errorf("Failed to fetch row: %v", err)
			}
			words = append(words, word)
		}
		rows.Close()
	}
	return strings.Join(words, ""), err
}
