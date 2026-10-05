package infraestructure

import (
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/ediciones/domain"
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

func (mysql *MySQL) SaveEdicion(tipoEdicion string, nombreEdicion string) (*domain.Edicion, error) {
	query := "INSERT INTO ediciones (tipo_edicion, nombre_edicion) VALUES (?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, tipoEdicion, nombreEdicion)
	if err != nil {
		return nil, fmt.Errorf("error al guardar edición: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID insertado: %v", err)
	}

	edicion := domain.NewEdicion(tipoEdicion, nombreEdicion)
	edicion.ID = int32(id)

	log.Printf("[MySQL] - Edición creada correctamente: ID:%d Nombre:%s", id, nombreEdicion)
	return edicion, nil
}

func (mysql *MySQL) GetAll() ([]domain.Edicion, error) {
	query := "SELECT id_edicion, tipo_edicion, nombre_edicion, fecha_inicio FROM ediciones ORDER BY id_edicion DESC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var ediciones []domain.Edicion

	for rows.Next() {
		var e domain.Edicion
		if err := rows.Scan(&e.ID, &e.TipoEdicion, &e.NombreEdicion, &e.FechaInicio); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %v", err)
		}
		ediciones = append(ediciones, e)
	}

	return ediciones, nil
}

func (mysql *MySQL) GetEdicionByID(id int32) (*domain.Edicion, error) {
	query := "SELECT id_edicion, tipo_edicion, nombre_edicion, fecha_inicio FROM ediciones WHERE id_edicion = ?"
	row := mysql.conn.FetchRow(query, id)

	var e domain.Edicion
	err := row.Scan(&e.ID, &e.TipoEdicion, &e.NombreEdicion, &e.FechaInicio)
	if err != nil {
		return nil, fmt.Errorf("edición no encontrada: %v", err)
	}

	return &e, nil
}

func (mysql *MySQL) UpdateEdicion(id int32, tipoEdicion string, nombreEdicion string) error {
	query := "UPDATE ediciones SET tipo_edicion = ?, nombre_edicion = ? WHERE id_edicion = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, tipoEdicion, nombreEdicion, id)
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