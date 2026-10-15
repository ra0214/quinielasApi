package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
	"time"
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
	query := "INSERT INTO quinielas (id_edicion, nombre_variante, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta, sujeto_cancelacion, estado_fallo, monto_meta_anterior, fecha_meta_anterior) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var idCliente interface{}
	if q.IDCliente != nil {
		idCliente = *q.IDCliente
	} else {
		idCliente = nil
	}
	result, err := mysql.conn.ExecutePreparedQuery(query, q.IDEdicion, q.NombreVariante, q.MontoMeta, q.FechaLimite, q.Estado, q.TipoQuiniela, idCliente, q.MontoPagado, q.Analista, q.TipoPago, q.CanalVenta, q.SujetoCancelacion, q.EstadoFallo, nullableMoney(q.MontoMetaAnterior), nullableTime(q.FechaMetaAnterior))
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
	var sujetoCancelacion sql.NullInt32
	var estadoFallo sql.NullInt32
	var montoMetaAnterior sql.NullString
	var fechaMetaAnterior sql.NullTime
	if err := scanner.Scan(&q.ID, &q.IDEdicion, &q.NombreVariante, &q.MontoMeta, &q.FechaLimite, &q.Estado, &tipoQ, &idC, &q.MontoPagado, &analista, &tipoPago, &canalVenta, &sujetoCancelacion, &estadoFallo, &montoMetaAnterior, &fechaMetaAnterior); err != nil {
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
	q.SujetoCancelacion = sujetoCancelacion.Valid && sujetoCancelacion.Int32 != 0
	q.EstadoFallo = int(estadoFallo.Int32)
	if montoMetaAnterior.Valid {
		if m, err := money.Parse(montoMetaAnterior.String); err == nil {
			q.MontoMetaAnterior = &m
		}
	}
	if fechaMetaAnterior.Valid {
		t := fechaMetaAnterior.Time
		q.FechaMetaAnterior = &t
	}
	return &q, nil
}

func nullableMoney(m *money.Money) interface{} {
	if m == nil {
		return nil
	}
	return *m
}

func nullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}

func (mysql *MySQL) GetAll() ([]domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta, sujeto_cancelacion, estado_fallo, monto_meta_anterior, fecha_meta_anterior FROM quinielas ORDER BY id_quiniela DESC"
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
	query := "SELECT id_quiniela, id_edicion, nombre_variante, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta, sujeto_cancelacion, estado_fallo, monto_meta_anterior, fecha_meta_anterior FROM quinielas WHERE id_quiniela = ?"
	row := mysql.conn.FetchRow(query, id)

	q, err := scanQuiniela(row)
	if err != nil {
		return nil, fmt.Errorf("quiniela no encontrada: %v", err)
	}

	return q, nil
}

func (mysql *MySQL) GetByEdicionID(idEdicion int32) ([]domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta, sujeto_cancelacion, estado_fallo, monto_meta_anterior, fecha_meta_anterior FROM quinielas WHERE id_edicion = ?"
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

// GetVencidasDesarrollo devuelve las quinielas DESARROLLO en juego cuya fecha
// límite ya pasó, para aplicarles la consecuencia del vencimiento.
func (mysql *MySQL) GetVencidasDesarrollo(ahora time.Time) ([]domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, monto_meta, fecha_limite, estado, tipo_quiniela, id_cliente, monto_pagado, analista, tipo_pago, canal_venta, sujeto_cancelacion, estado_fallo, monto_meta_anterior, fecha_meta_anterior FROM quinielas WHERE estado = 'EN_JUEGO' AND tipo_quiniela = 'DESARROLLO' AND fecha_limite <= ?"
	rows, err := mysql.conn.FetchRows(query, ahora)
	if err != nil {
		return nil, fmt.Errorf("error al obtener quinielas vencidas: %v", err)
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

func (mysql *MySQL) UpdateQuiniela(q *domain.Quiniela) error {
	query := "UPDATE quinielas SET nombre_variante = ?, monto_meta = ?, fecha_limite = ?, estado = ?, tipo_quiniela = ?, id_cliente = ?, monto_pagado = ?, analista = ?, tipo_pago = ?, canal_venta = ?, sujeto_cancelacion = ?, estado_fallo = ?, monto_meta_anterior = ?, fecha_meta_anterior = ? WHERE id_quiniela = ?"
	var idCliente interface{}
	if q.IDCliente != nil {
		idCliente = *q.IDCliente
	} else {
		idCliente = nil
	}
	result, err := mysql.conn.ExecutePreparedQuery(query, q.NombreVariante, q.MontoMeta, q.FechaLimite, q.Estado, q.TipoQuiniela, idCliente, q.MontoPagado, q.Analista, q.TipoPago, q.CanalVenta, q.SujetoCancelacion, q.EstadoFallo, nullableMoney(q.MontoMetaAnterior), nullableTime(q.FechaMetaAnterior), q.ID)
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