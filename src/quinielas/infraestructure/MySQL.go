package infraestructure

import (
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

func (mysql *MySQL) SaveQuiniela(idEdicion int32, nombreVariante string, precio money.Money, montoMeta money.Money, fechaLimite time.Time) (*domain.Quiniela, error) {
	query := "INSERT INTO quinielas (id_edicion, nombre_variante, precio, monto_meta, fecha_limite) VALUES (?, ?, ?, ?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, idEdicion, nombreVariante, precio, montoMeta, fechaLimite)
	if err != nil {
		return nil, fmt.Errorf("error al guardar quiniela: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID insertado: %v", err)
	}

	quiniela := domain.NewQuiniela(idEdicion, nombreVariante, precio, montoMeta, fechaLimite)
	quiniela.ID = int32(id)

	log.Printf("[MySQL] - Quiniela creada correctamente: ID:%d Variante:%s", id, nombreVariante)
	return quiniela, nil
}

func (mysql *MySQL) GetAll() ([]domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, precio, monto_meta, fecha_limite, estado FROM quinielas ORDER BY id_quiniela DESC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var quinielas []domain.Quiniela

	for rows.Next() {
		var q domain.Quiniela
		if err := rows.Scan(&q.ID, &q.IDEdicion, &q.NombreVariante, &q.Precio, &q.MontoMeta, &q.FechaLimite, &q.Estado); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		quinielas = append(quinielas, q)
	}

	return quinielas, nil
}

func (mysql *MySQL) GetQuinielaByID(id int32) (*domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, precio, monto_meta, fecha_limite, estado FROM quinielas WHERE id_quiniela = ?"
	row := mysql.conn.FetchRow(query, id)

	var q domain.Quiniela
	err := row.Scan(&q.ID, &q.IDEdicion, &q.NombreVariante, &q.Precio, &q.MontoMeta, &q.FechaLimite, &q.Estado)
	if err != nil {
		return nil, fmt.Errorf("quiniela no encontrada: %v", err)
	}

	return &q, nil
}

func (mysql *MySQL) GetByEdicionID(idEdicion int32) ([]domain.Quiniela, error) {
	query := "SELECT id_quiniela, id_edicion, nombre_variante, precio, monto_meta, fecha_limite, estado FROM quinielas WHERE id_edicion = ?"
	rows, err := mysql.conn.FetchRows(query, idEdicion)
	if err != nil {
		return nil, fmt.Errorf("error al obtener quinielas de la edición: %v", err)
	}
	defer rows.Close()

	var quinielas []domain.Quiniela

	for rows.Next() {
		var q domain.Quiniela
		if err := rows.Scan(&q.ID, &q.IDEdicion, &q.NombreVariante, &q.Precio, &q.MontoMeta, &q.FechaLimite, &q.Estado); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		quinielas = append(quinielas, q)
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

func (mysql *MySQL) UpdateQuiniela(id int32, nombreVariante string, precio money.Money, montoMeta money.Money, fechaLimite time.Time, estado string) error {
	query := "UPDATE quinielas SET nombre_variante = ?, precio = ?, monto_meta = ?, fecha_limite = ?, estado = ? WHERE id_quiniela = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, nombreVariante, precio, montoMeta, fechaLimite, estado, id)
	if err != nil {
		return fmt.Errorf("error al actualizar quiniela: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Quiniela actualizada correctamente: ID:%d", id)
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