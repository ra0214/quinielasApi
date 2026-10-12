package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/quinielas/domain"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IQuiniela = (*MySQL)(nil)

func NewMySQL() domain.IQuiniela {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveQuiniela(q *domain.Quiniela) (*domain.Quiniela, error) {
	query := "INSERT INTO quinielas (id_edicion, nombre_variante, precio, monto_meta, fecha_limite, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var idCliente interface{}
	if q.IDCliente != nil {
		idCliente = *q.IDCliente
	} else {
		idCliente = nil
	}
	result, err := mysql.conn.ExecutePreparedQuery(query, q.IDEdicion, q.NombreVariante, q.Precio, q.MontoMeta, q.FechaLimite, q.TipoQuiniela, idCliente, q.MontoPagado, q.Analista, q.TipoPago, q.CanalVenta)
	if err != nil {
		return nil, fmt.Errorf("error al guardar quiniela: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID insertado: %v", err)
	}

	q.ID = int32(id)

	log.Printf("[MySQL] - Quiniela creada correctamente: ID:%d Variante:%s", id, q.NombreVariante)
	return q, nil
}

// scanQuiniela lee una fila completa de `quinielas` tolerando NULL en las
// columnas nuevas (filas existentes creadas antes de la migración).
func scanQuiniela(scanner interface{ Scan(dest ...interface{}) error }) (*domain.Quiniela, error) {
	var q domain.Quiniela
	var idC sql.NullInt32
	var tipoQ sql.NullString
	var analista sql.NullString
	var tipoPago sql.NullString
	var canalVenta sql.NullString
	if err := scanner.Scan(&q.ID, &q.IDEdicion, &q.NombreVariante, &q.Precio, &q.MontoMeta, &q.FechaLimite, &q.Estado, &tipoQ, &idC, &q.MontoPagado, &analista, &tipoPago, &canalVenta); err != nil {
		return nil, err
	}
	if idC.Valid {
		v := idC.Int32
		q.IDCliente = &v
	}
	q.TipoQuiniela = tipoQ.String
	q.Analista = analista.String
	q.TipoPago = tipoPago.String
	q.CanalVenta = canalVenta.String
	return &q, nil
}

func (mysql *MySQL) GetAll() ([]domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, precio, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta FROM quinielas ORDER BY id_quiniela DESC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var quinielas []domain.Quiniela

	for rows.Next() {
		q, err := scanQuiniela(rows)
		if err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		quinielas = append(quinielas, *q)
	}

	return quinielas, nil
}

func (mysql *MySQL) GetQuinielaByID(id int32) (*domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, precio, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta FROM quinielas WHERE id_quiniela = ?"
	row := mysql.conn.FetchRow(query, id)

	q, err := scanQuiniela(row)
	if err != nil {
		return nil, fmt.Errorf("quiniela no encontrada: %v", err)
	}

	return q, nil
}

func (mysql *MySQL) GetByEdicionID(idEdicion int32) ([]domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, precio, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta FROM quinielas WHERE id_edicion = ?"
	rows, err := mysql.conn.FetchRows(query, idEdicion)
	if err != nil {
		return nil, fmt.Errorf("error al obtener quinielas de la edición: %v", err)
	}
	defer rows.Close()

	var quinielas []domain.Quiniela

	for rows.Next() {
		q, err := scanQuiniela(rows)
		if err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		quinielas = append(quinielas, *q)
	}

	return quinielas, nil
}

func (mysql *MySQL) CountByEdicionID(idEdicion int32) (int, error) {
	query := "SELECT COUNT(*) FROM quinielas WHERE id_edicion = ?"
	row := mysql.conn.FetchRow(query, idEdicion)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar quinielas por edición: %v", err)
	}

	return count, nil
}

func (mysql *MySQL) UpdateQuiniela(q *domain.Quiniela) error {
	query := "UPDATE quinielas SET nombre_variante = ?, precio = ?, monto_meta = ?, fecha_limite = ?, estado = ?, tipo_quiniela = ?, id_cliente = ?, monto_pagado = ?, analista = ?, tipo_pago = ?, canal_venta = ? WHERE id_quiniela = ?"
	var idCliente interface{}
	if q.IDCliente != nil {
		idCliente = *q.IDCliente
	} else {
		idCliente = nil
	}
	result, err := mysql.conn.ExecutePreparedQuery(query, q.NombreVariante, q.Precio, q.MontoMeta, q.FechaLimite, q.Estado, q.TipoQuiniela, idCliente, q.MontoPagado, q.Analista, q.TipoPago, q.CanalVenta, q.ID)
	if err != nil {
		return fmt.Errorf("error al actualizar quiniela: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Quiniela actualizada correctamente: ID:%d", q.ID)
	}
	return nil
}

func (mysql *MySQL) DeleteQuiniela(id int32) error {
	query := "DELETE FROM quinielas WHERE id_quiniela = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar la quiniela: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Quiniela eliminada correctamente: ID: %d", id)
	}
	return nil
}