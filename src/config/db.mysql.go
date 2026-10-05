package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Conn_MySQL struct {
	DB  *sql.DB
	Err string
}

func GetDBPool() *Conn_MySQL {
	error := ""
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSchema := os.Getenv("DB_SCHEMA")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPass, dbHost, dbSchema)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		error = fmt.Sprintf("error al abrir la base de datos: %v", err)
	}

	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		db.Close()
		error = fmt.Sprintf("error al verificar la conexión a la base de datos: %v", err)
	}

	return &Conn_MySQL{DB: db, Err: error}
}

func GetDBConnection() (*sql.DB, error) {
	conn := GetDBPool()
	if conn.Err != "" {
		return nil, fmt.Errorf("Error en la conexión a la base de datos: %s", conn.Err)
	}
	return conn.DB, nil
}

func (conn *Conn_MySQL) ExecutePreparedQuery(query string, values ...interface{}) (sql.Result, error) {
	stmt, err := conn.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error al preparar la consulta: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta preparada: %v", err)
	}

	return result, nil
}

func (conn *Conn_MySQL) FetchRows(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %v", err)
	}

	return rows, nil
}

func (conn *Conn_MySQL) FetchRow(query string, values ...interface{}) *sql.Row {
	return conn.DB.QueryRow(query, values...)
}
