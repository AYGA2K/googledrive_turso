package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
)

type Database struct {
	*sql.DB
	ctx context.Context
}

func ConnectDB() Database {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Generate encoded token and send it as response.
	dbUrl := os.Getenv("dbUrl")
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}
	return Database{
		db,
		context.Background(),
	}
}

// exec executes a SQL statement and returns the result.
func (db Database) exec(sql string, args ...interface{}) (sql.Result, error) {
	res, err := db.ExecContext(db.ctx, sql, args...)
	return res, err
}

// query executes a SQL query and returns the rows.
func (db Database) query(sql string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.QueryContext(db.ctx, sql, args...)
	return rows, err
}

type Table struct {
	Db   Database
	Name string
}

// CreateTable creates a new table with the specified column names and types.
func (db Database) CreateTable(tableName string, columns map[string]string) (Table, error) {
	columnDefinitions := make([]string, 0, len(columns))

	for col, colType := range columns {
		columnDefinitions = append(columnDefinitions, fmt.Sprintf("%s %s", col, colType))
	}

	query := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columnDefinitions, ", "))
	_, err := db.exec(query)
	if err != nil {
		return Table{}, err
	}
	return Table{
		db,
		tableName,
	}, err
}

func (t Table) Create(data map[string]interface{}) (int64, error) {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for col, val := range data {
		columns = append(columns, col)
		values = append(values, val)
	}

	query := fmt.Sprintf("INSERT INTO %s VALUES (%s)", t.Name, strings.Join(columns, ","))
	result, err := t.Db.exec(query, values...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
