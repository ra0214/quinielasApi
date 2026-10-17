package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"quinielas/src/config"
	"quinielas/src/mantenimiento/domain"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IPurga = (*MySQL)(nil)

func NewMySQL() domain.IPurga {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

// PurgarConcursosAnteriores elimina en una sola transaccion los concursos
// anteriores a "limite" y todo su contenido. El orden respeta las llaves
// foraneas: primero los hijos y al final las ediciones.
func (mysql *MySQL) PurgarConcursosAnteriores(limite time.Time) (int, error) {
	tx, err := mysql.conn.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("error al iniciar la transaccion de purga: %w", err)
	}

	ids, err := idsConcursosAnteriores(tx, limite)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if len(ids) == 0 {
		_ = tx.Rollback()
		return 0, nil
	}

	if err := borrarConcursos(tx, ids); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("error al confirmar la purga: %w", err)
	}
	return len(ids), nil
}

// idsConcursosAnteriores devuelve las ediciones cuyo cierre (o inicio, si no
// hay cierre registrado) es anterior al limite indicado.
func idsConcursosAnteriores(tx *sql.Tx, limite time.Time) ([]int64, error) {
	rows, err := tx.Query(
		"SELECT id_edicion FROM ediciones WHERE COALESCE(fecha_limite, fecha_inicio) < ?",
		limite,
	)
	if err != nil {
		return nil, fmt.Errorf("error al buscar concursos antiguos: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error al leer id de edicion: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar las ediciones: %w", err)
	}
	return ids, nil
}

// borrarConcursos elimina, por cada edicion, sus repartos, premios,
// movimientos, aportes y quinielas, y al final la propia edicion.
func borrarConcursos(tx *sql.Tx, ids []int64) error {
	// Sentencias en orden seguro para las llaves foraneas (los hijos primero).
	consultas := []string{
		`DELETE prp FROM premio_reparto_participantes prp
			JOIN premio_repartos pr ON pr.id_reparto = prp.id_reparto
			JOIN quinielas q ON q.id_quiniela = pr.id_quiniela
			WHERE q.id_edicion = ?`,
		`DELETE pr FROM premio_repartos pr
			JOIN quinielas q ON q.id_quiniela = pr.id_quiniela
			WHERE q.id_edicion = ?`,
		`DELETE p FROM premios p
			JOIN quinielas q ON q.id_quiniela = p.id_quiniela
			WHERE q.id_edicion = ?`,
		`DELETE m FROM movimientos m
			JOIN quinielas q ON q.id_quiniela = m.id_quiniela
			WHERE q.id_edicion = ?`,
		`DELETE a FROM aportes a
			JOIN quinielas q ON q.id_quiniela = a.id_quiniela
			WHERE q.id_edicion = ?`,
		`DELETE FROM quinielas WHERE id_edicion = ?`,
	}

	for _, id := range ids {
		for _, query := range consultas {
			if _, err := tx.Exec(query, id); err != nil {
				return fmt.Errorf("error al purgar el concurso #%d: %w", id, err)
			}
		}
	}

	stmt, err := tx.Prepare("DELETE FROM ediciones WHERE id_edicion = ?")
	if err != nil {
		return fmt.Errorf("error al preparar el borrado de ediciones: %w", err)
	}
	defer stmt.Close()

	for _, id := range ids {
		if _, err := stmt.Exec(id); err != nil {
			return fmt.Errorf("error al eliminar el concurso #%d: %w", id, err)
		}
	}
	return nil
}
