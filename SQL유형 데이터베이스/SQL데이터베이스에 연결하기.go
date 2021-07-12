package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	pgURL := os.Getenv("PGURL")
	if pgURL == "" {
		log.Fatal("PGURL empty")
	}
	//데이터베이스값을 연다
	// database/sql을 위한 postgres 드라이버를 지정한다.
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//데이터베이스에 성공적으로 연결됐는지 확인하기 위해
	//Ping메소드를 사용할 수 있음
	if err := db.Ping(); err != nil {
		log.Fatal(err)

	}

}
