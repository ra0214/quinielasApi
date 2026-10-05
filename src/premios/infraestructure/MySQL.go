package infraestructure

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/premios/domain"
	"quinielas/src/shared/money"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IPremio = (*MySQL)(nil)

func NewMySQL() domain.IPremio {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SavePremio(idQuiniela int32, montoBruto money.Money) (*domain.Premio, error) {
	nuevoPremio := domain.NewPremio(idQuiniela, montoBruto)

	query := `INSERT INTO premios (id_quiniela, monto_bruto, porcentaje_retencion, monto_neto) 
	          VALUES (?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE monto_bruto = ?, porcentaje_retencion = ?, monto_neto = ?`

	res, err := mysql.conn.ExecutePreparedQuery(
		query,
		idQuiniela,
		nuevoPremio.MontoBruto,
		nuevoPremio.PorcentajeRetencion,
		nuevoPremio.MontoNeto,
		nuevoPremio.MontoBruto,
		nuevoPremio.PorcentajeRetencion,
		nuevoPremio.MontoNeto,
	)
	if err != nil {
		return nil, fmt.Errorf("error al registrar premio: %v", err)
	}

	id, _ := res.LastInsertId()
	nuevoPremio.IDPremio = int32(id)

	log.Printf("[MySQL] - Premio registrado correctamente: Quiniela:%d Neto:%s", idQuiniela, nuevoPremio.MontoNeto.String())
	return nuevoPremio, nil
}

func (mysql *MySQL) DeletePremio(idPremio int32) error {
	query := "DELETE FROM premios WHERE id_premio = ?"
	_, err := mysql.conn.ExecutePreparedQuery(query, idPremio)
	if err != nil {
		return fmt.Errorf("error al eliminar premio: %v", err)
	}
	return nil
}

func (mysql *MySQL) GetPremioByQuinielaID(idQuiniela int32) (*domain.Premio, error) {
	query := "SELECT id_premio, id_quiniela, monto_bruto, porcentaje_retencion, monto_neto, fecha_registro FROM premios WHERE id_quiniela = ?"
	row := mysql.conn.FetchRow(query, idQuiniela)

	var p domain.Premio
	err := row.Scan(&p.IDPremio, &p.IDQuiniela, &p.MontoBruto, &p.PorcentajeRetencion, &p.MontoNeto, &p.FechaRegistro)
	if err != nil {
		// sql.ErrNoRows se devuelve sin envolver para que el caso de uso lo
		// distinga de un fallo real de conexion.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("error al consultar el premio: %w", err)
	}

	return &p, nil
}

func (mysql *MySQL) GetAllPremios() ([]domain.Premio, error) {
	query := "SELECT id_premio, id_quiniela, monto_bruto, porcentaje_retencion, monto_neto, fecha_registro FROM premios ORDER BY fecha_registro DESC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener premios: %v", err)
	}
	defer rows.Close()

	var lista []domain.Premio
	for rows.Next() {
		var p domain.Premio
		if err := rows.Scan(&p.IDPremio, &p.IDQuiniela, &p.MontoBruto, &p.PorcentajeRetencion, &p.MontoNeto, &p.FechaRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear fila de premio: %v", err)
		}
		lista = append(lista, p)
	}

	return lista, nil
}