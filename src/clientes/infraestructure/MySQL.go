package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/clientes/domain"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.ICliente = (*MySQL)(nil)

func NewMySQL() domain.ICliente {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveCliente(nombre string, telefono string) (*domain.Cliente, error) {
	query := "INSERT INTO clientes (nombre, telefono) VALUES (?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, nombre, telefono)
	if err != nil {
		return nil, fmt.Errorf("error al guardar cliente: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID insertado: %v", err)
	}

	cliente := domain.NewCliente(nombre, telefono)
	cliente.ID = int32(id)

	log.Printf("[MySQL] - Cliente creado correctamente: ID:%d Nombre:%s", id, nombre)
	return cliente, nil
}

func (mysql *MySQL) GetAll() ([]domain.Cliente, error) {
	query := "SELECT id_cliente, nombre, telefono, fecha_registro FROM clientes ORDER BY nombre ASC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var clientes []domain.Cliente

	for rows.Next() {
		var c domain.Cliente
		var telefono sql.NullString
		if err := rows.Scan(&c.ID, &c.Nombre, &telefono, &c.FechaRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		if telefono.Valid {
			c.Telefono = telefono.String
		}
		clientes = append(clientes, c)
	}

	return clientes, nil
}

func (mysql *MySQL) GetClienteByID(id int32) (*domain.Cliente, error) {
	query := "SELECT id_cliente, nombre, telefono, fecha_registro FROM clientes WHERE id_cliente = ?"
	row := mysql.conn.FetchRow(query, id)

	var c domain.Cliente
	var telefono sql.NullString
	err := row.Scan(&c.ID, &c.Nombre, &telefono, &c.FechaRegistro)
	if err != nil {
		return nil, fmt.Errorf("cliente no encontrado: %v", err)
	}
	if telefono.Valid {
		c.Telefono = telefono.String
	}

	return &c, nil
}

func (mysql *MySQL) SearchClientesByName(searchQuery string) ([]domain.Cliente, error) {
	query := "SELECT id_cliente, nombre, telefono, fecha_registro FROM clientes WHERE nombre LIKE ? ORDER BY nombre ASC"
	rows, err := mysql.conn.FetchRows(query, "%"+searchQuery+"%")
	if err != nil {
		return nil, fmt.Errorf("error al buscar clientes: %v", err)
	}
	defer rows.Close()

	var clientes []domain.Cliente

	for rows.Next() {
		var c domain.Cliente
		var telefono sql.NullString
		if err := rows.Scan(&c.ID, &c.Nombre, &telefono, &c.FechaRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		if telefono.Valid {
			c.Telefono = telefono.String
		}
		clientes = append(clientes, c)
	}

	return clientes, nil
}

func (mysql *MySQL) UpdateCliente(id int32, nombre string, telefono string) error {
	query := "UPDATE clientes SET nombre = ?, telefono = ? WHERE id_cliente = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, nombre, telefono, id)
	if err != nil {
		return fmt.Errorf("error al actualizar cliente: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Cliente actualizado correctamente: ID:%d Nombre:%s", id, nombre)
	}
	return nil
}

func (mysql *MySQL) DeleteCliente(id int32) error {
	// Las tablas saldos, aportes y movimientos referencian a clientes con FK,
	// así que primero limpiamos esas dependencias dentro de una transacción.
	tx, err := mysql.conn.DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %v", err)
	}

	for _, tabla := range []string{"saldos", "movimientos", "aportes"} {
		query := "DELETE FROM " + tabla + " WHERE id_cliente = ?"
		if _, err := tx.Exec(query, id); err != nil {
			tx.Rollback()
			return fmt.Errorf("error al eliminar dependencias de %s: %v", tabla, err)
		}
	}

	query := "DELETE FROM clientes WHERE id_cliente = ?"
	result, err := tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error al eliminar el cliente: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Cliente eliminado correctamente: ID: %d", id)
	}
	return nil
}