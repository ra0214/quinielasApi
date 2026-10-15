package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/ediciones/domain"
	"time"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IEdicion = (*MySQL)(nil)

func NewMySQL() domain.IEdicion {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveEdicion(tipoEdicion string, nombreEdicion string, fechaLimite *time.Time) (*domain.Edicion, error) {
	query := "INSERT INTO ediciones (tipo_edicion, nombre_edicion, fecha_limite) VALUES (?, ?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, tipoEdicion, nombreEdicion, fechaLimite)
	if err != nil {
		return nil, fmt.Errorf("error al guardar edición: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID insertado: %v", err)
	}

	edicion := domain.NewEdicion(tipoEdicion, nombreEdicion, fechaLimite)
	edicion.ID = int32(id)

	log.Printf("[MySQL] - Edición creada correctamente: ID:%d Nombre:%s", id, nombreEdicion)
	return edicion, nil
}

func (mysql *MySQL) GetAll() ([]domain.Edicion, error) {
	query := "SELECT id_edicion, tipo_edicion, nombre_edicion, fecha_inicio, fecha_limite FROM ediciones ORDER BY id_edicion DESC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var ediciones []domain.Edicion

	for rows.Next() {
		var e domain.Edicion
		var fechaLimite sql.NullTime
		if err := rows.Scan(&e.ID, &e.TipoEdicion, &e.NombreEdicion, &e.FechaInicio, &fechaLimite); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		if fechaLimite.Valid {
			e.FechaLimite = &fechaLimite.Time
		}
		ediciones = append(ediciones, e)
	}

	return ediciones, nil
}

func (mysql *MySQL) GetEdicionByID(id int32) (*domain.Edicion, error) {
	query := "SELECT id_edicion, tipo_edicion, nombre_edicion, fecha_inicio, fecha_limite FROM ediciones WHERE id_edicion = ?"
	row := mysql.conn.FetchRow(query, id)

	var e domain.Edicion
	var fechaLimite sql.NullTime
	if err := row.Scan(&e.ID, &e.TipoEdicion, &e.NombreEdicion, &e.FechaInicio, &fechaLimite); err != nil {
		return nil, fmt.Errorf("error al obtener la edición: %v", err)
	}
	if fechaLimite.Valid {
		e.FechaLimite = &fechaLimite.Time
	}

	return &e, nil
}

func (mysql *MySQL) UpdateEdicion(id int32, tipoEdicion string, nombreEdicion string, fechaLimite *time.Time) error {
	query := "UPDATE ediciones SET tipo_edicion = ?, nombre_edicion = ?, fecha_limite = ? WHERE id_edicion = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, tipoEdicion, nombreEdicion, fechaLimite, id)
	if err != nil {
		return fmt.Errorf("error al actualizar edición: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Edición actualizada correctamente: ID:%d", id)
	}
	return nil
}

func (mysql *MySQL) DeleteEdicion(id int32) error {
	query := "DELETE FROM ediciones WHERE id_edicion = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar la edición: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Edición eliminada correctamente: ID: %d", id)
	}
	return nil
}
