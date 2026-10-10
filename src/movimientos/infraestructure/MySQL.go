package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/movimientos/domain"
	"quinielas/src/shared/money"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IMovimiento = (*MySQL)(nil)

func NewMySQL() domain.IMovimiento {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveMovimiento(idCliente int32, idQuiniela *int32, tipo string, monto money.Money, descripcion string) (*domain.Movimiento, error) {
	query := "INSERT INTO movimientos (id_cliente, id_quiniela, tipo, monto, descripcion) VALUES (?, ?, ?, ?, ?)"
	
	var res sql.Result
	var err error

	if idQuiniela != nil {
		res, err = mysql.conn.ExecutePreparedQuery(query, idCliente, *idQuiniela, tipo, monto, descripcion)
	} else {
		res, err = mysql.conn.ExecutePreparedQuery(query, idCliente, nil, tipo, monto, descripcion)
	}

	if err != nil {
		return nil, fmt.Errorf("error al guardar movimiento: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID del movimiento: %v", err)
	}

	movimiento := domain.NewMovimiento(idCliente, idQuiniela, tipo, monto, descripcion)
	movimiento.IDMovimiento = int32(id)

	log.Printf("[MySQL] - Movimiento registrado correctamente: ID:%d Tipo:%s Cliente:%d", id, tipo, idCliente)
	return movimiento, nil
}

func (mysql *MySQL) UpdateMovimiento(m *domain.Movimiento) error {
	query := "UPDATE movimientos SET id_cliente = ?, id_quiniela = ?, tipo = ?, monto = ?, descripcion = ? WHERE id_movimiento = ?"

	var err error
	if m.IDQuiniela != nil {
		_, err = mysql.conn.ExecutePreparedQuery(query, m.IDCliente, *m.IDQuiniela, m.Tipo, m.Monto, m.Descripcion, m.IDMovimiento)
	} else {
		_, err = mysql.conn.ExecutePreparedQuery(query, m.IDCliente, nil, m.Tipo, m.Monto, m.Descripcion, m.IDMovimiento)
	}

	if err != nil {
		return fmt.Errorf("error al actualizar movimiento: %v", err)
	}

	log.Printf("[MySQL] - Movimiento actualizado: ID:%d Tipo:%s Cliente:%d", m.IDMovimiento, m.Tipo, m.IDCliente)
	return nil
}

func (mysql *MySQL) DeleteMovimiento(idMovimiento int32) error {
	query := "DELETE FROM movimientos WHERE id_movimiento = ?"
	_, err := mysql.conn.ExecutePreparedQuery(query, idMovimiento)
	if err != nil {
		return fmt.Errorf("error al eliminar movimiento: %v", err)
	}
	return nil
}

func (mysql *MySQL) GetMovimientoByID(idMovimiento int32) (*domain.Movimiento, error) {
	query := "SELECT id_movimiento, id_cliente, id_quiniela, tipo, monto, descripcion, fecha FROM movimientos WHERE id_movimiento = ?"
	row := mysql.conn.FetchRow(query, idMovimiento)

	var m domain.Movimiento
	var idQ sql.NullInt32
	if err := row.Scan(&m.IDMovimiento, &m.IDCliente, &idQ, &m.Tipo, &m.Monto, &m.Descripcion, &m.Fecha); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movimiento %d no encontrado", idMovimiento)
		}
		return nil, fmt.Errorf("error al obtener movimiento: %v", err)
	}
	if idQ.Valid {
		val := idQ.Int32
		m.IDQuiniela = &val
	}

	return &m, nil
}

func (mysql *MySQL) GetMovimientosByClienteID(idCliente int32) ([]domain.Movimiento, error) {
	query := "SELECT id_movimiento, id_cliente, id_quiniela, tipo, monto, descripcion, fecha FROM movimientos WHERE id_cliente = ? ORDER BY fecha DESC"
	rows, err := mysql.conn.FetchRows(query, idCliente)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos del cliente: %v", err)
	}
	defer rows.Close()

	var lista []domain.Movimiento
	for rows.Next() {
		var m domain.Movimiento
		var idQ sql.NullInt32
		if err := rows.Scan(&m.IDMovimiento, &m.IDCliente, &idQ, &m.Tipo, &m.Monto, &m.Descripcion, &m.Fecha); err != nil {
			return nil, fmt.Errorf("error al escanear fila de movimiento: %v", err)
		}
		if idQ.Valid {
			val := idQ.Int32
			m.IDQuiniela = &val
		}
		lista = append(lista, m)
	}

	return lista, nil
}

func (mysql *MySQL) GetAllMovimientos() ([]domain.Movimiento, error) {
	query := "SELECT id_movimiento, id_cliente, id_quiniela, tipo, monto, descripcion, fecha FROM movimientos ORDER BY fecha DESC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos: %v", err)
	}
	defer rows.Close()

	var lista []domain.Movimiento
	for rows.Next() {
		var m domain.Movimiento
		var idQ sql.NullInt32
		if err := rows.Scan(&m.IDMovimiento, &m.IDCliente, &idQ, &m.Tipo, &m.Monto, &m.Descripcion, &m.Fecha); err != nil {
			return nil, fmt.Errorf("error al escanear fila de movimiento: %v", err)
		}
		if idQ.Valid {
			val := idQ.Int32
			m.IDQuiniela = &val
		}
		lista = append(lista, m)
	}

	return lista, nil
}