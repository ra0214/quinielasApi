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

func (mysql *MySQL) SavePremio(idQuiniela int32, montoBruto money.Money, premioPorBoleto money.Money) (*domain.Premio, error) {
	var nuevoPremio *domain.Premio
	if premioPorBoleto.GreaterThan(money.Zero()) {
		nuevoPremio = domain.NewPremioCoperacha(idQuiniela, montoBruto, premioPorBoleto)
	} else {
		nuevoPremio = domain.NewPremio(idQuiniela, montoBruto)
	}

	query := `INSERT INTO premios (id_quiniela, monto_bruto, porcentaje_retencion, monto_neto, premio_por_boleto) 
	          VALUES (?, ?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE monto_bruto = ?, porcentaje_retencion = ?, monto_neto = ?, premio_por_boleto = ?`

	res, err := mysql.conn.ExecutePreparedQuery(
		query,
		idQuiniela,
		nuevoPremio.MontoBruto,
		nuevoPremio.PorcentajeRetencion,
		nuevoPremio.MontoNeto,
		nuevoPremio.PremioPorBoleto,
		nuevoPremio.MontoBruto,
		nuevoPremio.PorcentajeRetencion,
		nuevoPremio.MontoNeto,
		nuevoPremio.PremioPorBoleto,
	)
	if err != nil {
		return nil, fmt.Errorf("error al registrar premio: %v", err)
	}

	id, _ := res.LastInsertId()
	nuevoPremio.IDPremio = int32(id)

	log.Printf("[MySQL] - Premio registrado correctamente: Quiniela:%d Neto:%s", idQuiniela, nuevoPremio.MontoNeto.String())
	return nuevoPremio, nil
}

func (mysql *MySQL) UpdatePremio(idPremio int32, montoBruto money.Money, premioPorBoleto money.Money) (*domain.Premio, error) {
	var premio *domain.Premio
	if premioPorBoleto.GreaterThan(money.Zero()) {
		// Necesitamos la quiniela para recalcular el total nominal. Si no se
		// encuentra seguimos con el monto capturado como total.
		premio = domain.NewPremioCoperacha(0, montoBruto, premioPorBoleto)
	} else {
		premio = domain.NewPremio(0, montoBruto)
	}
	premio.IDPremio = idPremio

	query := "UPDATE premios SET monto_bruto = ?, porcentaje_retencion = ?, monto_neto = ?, premio_por_boleto = ? WHERE id_premio = ?"
	res, err := mysql.conn.ExecutePreparedQuery(
		query,
		premio.MontoBruto,
		premio.PorcentajeRetencion,
		premio.MontoNeto,
		premio.PremioPorBoleto,
		idPremio,
	)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar premio: %v", err)
	}

	// Si id_premio no existe el UPDATE no afecta filas: lo tratamos como no
	// encontrado para que el cliente reciba un 404 claro.
	if filas, err := res.RowsAffected(); err == nil && filas == 0 {
		return nil, sql.ErrNoRows
	}

	log.Printf("[MySQL] - Premio actualizado: ID:%d Neto:%s", idPremio, premio.MontoNeto.String())
	return premio, nil
}

func (mysql *MySQL) DeletePremio(idPremio int32) error {
	query := "DELETE FROM premios WHERE id_premio = ?"
	_, err := mysql.conn.ExecutePreparedQuery(query, idPremio)
	if err != nil {
		return fmt.Errorf("error al eliminar premio: %v", err)
	}
	return nil
}

func (mysql *MySQL) GetPremioByID(idPremio int32) (*domain.Premio, error) {
	query := "SELECT id_premio, id_quiniela, monto_bruto, porcentaje_retencion, monto_neto, premio_por_boleto, fecha_registro FROM premios WHERE id_premio = ?"
	row := mysql.conn.FetchRow(query, idPremio)

	var p domain.Premio
	err := row.Scan(&p.IDPremio, &p.IDQuiniela, &p.MontoBruto, &p.PorcentajeRetencion, &p.MontoNeto, &p.PremioPorBoleto, &p.FechaRegistro)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("error al consultar el premio: %w", err)
	}

	return &p, nil
}

func (mysql *MySQL) GetPremioByQuinielaID(idQuiniela int32) (*domain.Premio, error) {
	query := "SELECT id_premio, id_quiniela, monto_bruto, porcentaje_retencion, monto_neto, premio_por_boleto, fecha_registro FROM premios WHERE id_quiniela = ?"
	row := mysql.conn.FetchRow(query, idQuiniela)

	var p domain.Premio
	err := row.Scan(&p.IDPremio, &p.IDQuiniela, &p.MontoBruto, &p.PorcentajeRetencion, &p.MontoNeto, &p.PremioPorBoleto, &p.FechaRegistro)
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
	query := "SELECT id_premio, id_quiniela, monto_bruto, porcentaje_retencion, monto_neto, premio_por_boleto, fecha_registro FROM premios ORDER BY fecha_registro DESC"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener premios: %v", err)
	}
	defer rows.Close()

	var lista []domain.Premio
	for rows.Next() {
		var p domain.Premio
		if err := rows.Scan(&p.IDPremio, &p.IDQuiniela, &p.MontoBruto, &p.PorcentajeRetencion, &p.MontoNeto, &p.PremioPorBoleto, &p.FechaRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear fila de premio: %v", err)
		}
		lista = append(lista, p)
	}

	return lista, nil
}