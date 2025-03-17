package testwork

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"time"
)

func insert(db *sql.DB) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomString := make([]byte, 1000)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	data := make([]string, 10)
	for i := range data {
		data[i] = string(randomString)
	}

	for j := 0; j < 1600; j++ {
		stmt, err := db.Prepare("INSERT INTO test_conn_s (field1, field2, field3, field4, field5, field6, " +
			"field7, field8, field9, field10) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalf("Failed to prepare SQL statement: %v", err)
		}
		for i := 0; i < 500; i++ {
			_, err = stmt.Exec(data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9])
			if err != nil {
				log.Fatalf("Failed to execute SQL statement: %v", err)
			}
		}
		stmt.Close()
		fmt.Println("Data inserted successfully!")
	}

}

func Select(db *sql.DB) {
	start := time.Now()
	rows, err := db.Query("SELECT *, CONCAT(field1, field3) AS f from (SELECT * FROM test_conn_s ORDER BY field1) ORDER BY f")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	//for rows.Next() {
	//	data := make([]string, 30)
	//	if err := rows.Scan(&data[0], &data[1], &data[2], &data[3], &data[4], &data[5], &data[6], &data[7], &data[8],
	//		&data[9], &data[10], &data[11], &data[12], &data[13], &data[14], &data[15], &data[16], &data[17],
	//		&data[18], &data[19], &data[20], &data[21], &data[22], &data[23], &data[24], &data[25], &data[26],
	//		&data[27], &data[28], &data[29]); err != nil {
	//		fmt.Println("Error scanning row:", err)
	//		continue
	//	}
	//}
	elapsed := time.Since(start)
	fmt.Println("The check took ", elapsed)
	defer rows.Close()
	fmt.Println("Select successfully!")
}

func Sleep(db *sql.DB) {
	_, err := db.Exec("SELECT SLEEP(?)", 10)
	if err != nil {
		log.Fatalf("Failed to execute SLEEP(): %v", err)
	}
}

func SelectMysql(db *sql.DB) {
	rows, err := db.Query("SELECT User from user")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
	}
	defer rows.Close()
	fmt.Println("Select successfully!")
}

func main() {
	db, err := sql.Open("mysql", "snan:19990928@tcp(localhost:3306)/mysql")
	if err != nil {
		log.Fatalf("Failsed to connect to database: %v", err)
	}
	//Select(db)
	//insert(db)
	//Sleep(db)
	SelectMysql(db)
}
