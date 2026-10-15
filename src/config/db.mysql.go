package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Conn_MySQL struct {
	DB  *sql.DB
	Err string
}

func GetDBPool() *Conn_MySQL {
	error := ""

	// Simplemente ignora el error si no existe el archivo .env
	_ = godotenv.Load(".env")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSchema := os.Getenv("DB_SCHEMA")

	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=5s", dbUser, dbPass, dbHost, dbPort, dbSchema)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		error = fmt.Sprintf("error al abrir la base de datos: %v", err)
		return &Conn_MySQL{Err: error}
	}

	db.SetMaxOpenConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
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
		return nil, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta preparada: %w", err)
	}

	return result, nil
}

func (conn *Conn_MySQL) FetchRows(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %w", err)
	}

	return rows, nil
}

func (conn *Conn_MySQL) FetchRow(query string, values ...interface{}) *sql.Row {
	return conn.DB.QueryRow(query, values...)
}

// AutoAlterQuinielas asegura que existan las columnas nuevas de la tabla
// quinielas. MySQL 8 no soporta "ADD COLUMN IF NOT EXISTS", asi que se
// comprueba information_schema antes de cada ALTER para que sea idempotente.
func (conn *Conn_MySQL) AutoAlterQuinielas() error {
	columnas := []struct {
		nombre string
		ddl    string
	}{
		{"tipo_quiniela", "ALTER TABLE quinielas ADD COLUMN tipo_quiniela VARCHAR(30) DEFAULT NULL"},
		{"id_cliente", "ALTER TABLE quinielas ADD COLUMN id_cliente INT DEFAULT NULL"},
		{"monto_pagado", "ALTER TABLE quinielas ADD COLUMN monto_pagado DECIMAL(12,2) NOT NULL DEFAULT '0.00'"},
		{"analista", "ALTER TABLE quinielas ADD COLUMN analista VARCHAR(100) DEFAULT NULL"},
		{"tipo_pago", "ALTER TABLE quinielas ADD COLUMN tipo_pago VARCHAR(30) DEFAULT NULL"},
		{"canal_venta", "ALTER TABLE quinielas ADD COLUMN canal_venta VARCHAR(20) DEFAULT NULL"},
	}
	for _, col := range columnas {
		existe, err := conn.existeColumna("quinielas", col.nombre)
		if err != nil {
			return err
		}
		if !existe {
			if _, err := conn.DB.Exec(col.ddl); err != nil {
				return fmt.Errorf("error al agregar columna %s de quinielas: %w", col.nombre, err)
			}
		}
	}

	// Indice y llave foranea sobre id_cliente (tolerantes a que ya existan).
	existeIdx, err := conn.existeIndice("quinielas", "idx_quinielas_id_cliente")
	if err != nil {
		return err
	}
	if !existeIdx {
		if _, err := conn.DB.Exec("ALTER TABLE quinielas ADD INDEX idx_quinielas_id_cliente (id_cliente)"); err != nil {
			return fmt.Errorf("error al crear indice idx_quinielas_id_cliente: %w", err)
		}
	}

	existeFK, err := conn.existeFK("quinielas", "fk_quinielas_cliente")
	if err != nil {
		return err
	}
	if !existeFK {
		if _, err := conn.DB.Exec("ALTER TABLE quinielas ADD CONSTRAINT fk_quinielas_cliente FOREIGN KEY (id_cliente) REFERENCES clientes(id_cliente)"); err != nil {
			return fmt.Errorf("error al crear fk_quinielas_cliente: %w", err)
		}
	}

	return nil
}

// AutoAlterEdiciones asegura que exista la columna fecha_limite (cierre del
// concurso) en la tabla ediciones, de forma idempotente.
func (conn *Conn_MySQL) AutoAlterEdiciones() error {
	existe, err := conn.existeColumna("ediciones", "fecha_limite")
	if err != nil {
		return err
	}
	if existe {
		return nil
	}
	if _, err := conn.DB.Exec("ALTER TABLE ediciones ADD COLUMN fecha_limite TIMESTAMP NULL DEFAULT NULL"); err != nil {
		return fmt.Errorf("error al agregar fecha_limite a ediciones: %w", err)
	}
	return nil
}

func (conn *Conn_MySQL) existeColumna(tabla, columna string) (bool, error) {
	var n int
	err := conn.DB.QueryRow(
		"SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?",
		tabla, columna,
	).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("error al consultar information_schema.COLUMNS: %w", err)
	}
	return n > 0, nil
}

func (conn *Conn_MySQL) existeIndice(tabla, indice string) (bool, error) {
	var n int
	err := conn.DB.QueryRow(
		"SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND INDEX_NAME = ?",
		tabla, indice,
	).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("error al consultar information_schema.STATISTICS: %w", err)
	}
	return n > 0, nil
}

func (conn *Conn_MySQL) existeFK(tabla, nombre string) (bool, error) {
	var n int
	err := conn.DB.QueryRow(
		"SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = ? AND CONSTRAINT_NAME = ? AND CONSTRAINT_TYPE = 'FOREIGN KEY'",
		tabla, nombre,
	).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("error al consultar information_schema.TABLE_CONSTRAINTS: %w", err)
	}
	return n > 0, nil
}
