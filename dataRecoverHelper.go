package agent

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DataRecoverHelper ...
type DataRecoverHelper interface {
	IsDataExist() bool
	Read(count int) []string
	Write(message string) bool
}

type dataRecoverHelper struct {
	lock     sync.Mutex
	filePath string
}

// NewDataRecoverHelper ...
func NewDataRecoverHelper(path string) DataRecoverHelper {
	return &dataRecoverHelper{
		filePath: path,
	}
}

func (helper *dataRecoverHelper) IsDataExist() bool {
	if _, err := os.Stat(helper.filePath); os.IsNotExist(err) {
		return false
	}
	helper.lock.Lock()
	db, err := sql.Open("sqlite3", helper.filePath)

	defer func() {
		helper.lock.Unlock()
		db.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return false
	}
	rows, err := db.Query("SELECT * FROM Data LIMIT 1")
	if err != nil {
		fmt.Println(err)
		return false
	}
	result := false
	if rows != nil {
		for rows.Next() {
			result = true
		}
	}
	return result
}

func (helper *dataRecoverHelper) Read(count int) []string {
	var messages []string
	var emptyMessages []string
	var ids []int
	if _, err := os.Stat(helper.filePath); os.IsNotExist(err) {
		return emptyMessages
	}

	helper.lock.Lock()
	db, err := sql.Open("sqlite3", helper.filePath)
	defer func() {
		helper.lock.Unlock()
		db.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return emptyMessages
	}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM Data LIMIT %d", count))
	if err != nil {
		fmt.Println(err)
		return emptyMessages
	}

	if rows == nil {
		return emptyMessages
	}

	for rows.Next() {
		var id int
		var message string
		err = rows.Scan(&id, &message)
		if err != nil {
			fmt.Println(err)
		}
		messages = append(messages, message)
		ids = append(ids, id)
	}
	var str string
	for _, value := range ids {
		str += strconv.Itoa(value) + ","
	}

	sql, err := db.Prepare("DELETE FROM Data WHERE id IN (" + str[:len(str)-1] + ")")
	if err != nil {
		fmt.Println(err)
		return emptyMessages
	}

	_, err = sql.Exec()
	if err != nil {
		fmt.Println(err)
		return emptyMessages
	}

	return messages
}

func (helper *dataRecoverHelper) Write(message string) bool {
	helper.lock.Lock()
	db, err := sql.Open("sqlite3", helper.filePath)

	defer func() {
		helper.lock.Unlock()
		db.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return false
	}
	sql, err := db.Prepare("CREATE TABLE IF NOT EXISTS Data (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, message TEXT NOT NULL)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = sql.Exec()
	if err != nil {
		fmt.Println(err)
		return false
	}
	sql, err = db.Prepare("INSERT INTO Data(message) VALUES(?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = sql.Exec(message)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
