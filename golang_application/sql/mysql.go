package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	zipkinsql "github.com/openzipkin-contrib/zipkin-go-sql"
	zipkin "github.com/openzipkin/zipkin-go"
	"log"
	"time"
)

var dbSource = "root:password@tcp(localhost:3306)/test_db"

func MySQLInit(tracer *zipkin.Tracer) error {
	driverName, err := zipkinsql.Register("mysql", tracer, zipkinsql.WithAllTraceOptions())
	if err != nil {
		log.Fatalf("unable to register zipkin driver: %v\n", err)
	}

	db, err := sql.Open(driverName, dbSource)
	defer db.Close()

	if err != nil {
		log.Println("sql open error : ", err.Error())
		return err
	}
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS fruit (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Println("create table error : ", err.Error())
		return err
	}

	_, err = db.Query("INSERT INTO fruit (name) VALUES ('apple');")
	if err != nil {
		return err
	}

	return nil
}

func GetFruits(tracer *zipkin.Tracer) string {
	driverName, err := zipkinsql.Register("mysql", tracer, zipkinsql.WithAllTraceOptions())
	if err != nil {
		log.Fatalf("unable to register zipkin driver: %v\n", err)
	}

	db, err := sql.Open(driverName, dbSource)
	defer db.Close()
	if err != nil {
		return err.Error()
	}

	rows, err := db.Query(`select * from fruit order by id desc limit 1`)
	if err != nil {
		return err.Error()
	}

	for rows.Next() {
		var id int
		var name string
		var createdAt time.Time
		rows.Scan(&id, &name, &createdAt)
		return fmt.Sprintf("id: %d, name: %s, created_at: %s", id, name, createdAt)
	}
	return "no data"
}
