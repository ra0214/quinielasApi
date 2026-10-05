package infraestructure

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"quinielas/src/premios/domain"
)

func (mysql *MySQL) ExistsReparto(idQuiniela int32) (bool, error) {
	query := "SELECT COUNT(*) FROM premio_repartos WHERE id_quiniela = ?"
	row := mysql.conn.FetchRow(query, idQuiniela)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("error al verificar si la quiniela ya fue repartida: %w", err)
	}
	return count > 0, nil
}

// SaveReparto congela el reparto en dos tablas: la cabecera con los totales y
// el detalle con la parte de cada cliente. Se escriben dentro de la misma
// transaccion, de modo que nunca queda un reparto con la mitad de los
// participantes guardado.
func (mysql *MySQL) SaveReparto(reparto *domain.RepartoPremio) (err error) {
	tx, err := mysql.conn.DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar la transaccion del reparto: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	queryHeader := `INSERT INTO premio_repartos
		(id_premio, id_quiniela, monto_bruto, porcentaje_retencion, monto_neto,
		 total_recaudado, monto_meta, porcentaje_meta_alcanzado, total_participantes, fecha_reparto)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	header, err := tx.Exec(queryHeader,
		reparto.IDPremio,
		reparto.IDQuiniela,
		reparto.MontoBruto,
		reparto.PorcentajeRetencion,
		reparto.MontoNeto,
		reparto.TotalRecaudado,
		reparto.MontoMeta,
		reparto.PorcentajeMetaAlcanzado,
		reparto.TotalParticipantes,
		reparto.FechaReparto,
	)
	if err != nil {
		return fmt.Errorf("error al guardar la cabecera del reparto: %w", err)
	}

	idReparto, err := header.LastInsertId()
	if err != nil {
		return fmt.Errorf("error al obtener el ID del reparto: %w", err)
	}

	queryDetalle := `INSERT INTO premio_reparto_participantes
		(id_reparto, id_cliente, nombre_cliente, monto_aporte, porcentaje_aporte,
		 porcentaje_sobre_meta, porcentaje_recompensa, monto_bruto_asignado,
		 monto_retencion, monto_neto_asignado)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, p := range reparto.Participantes {
		if _, err := tx.Exec(queryDetalle,
			idReparto,
			p.IDCliente,
			p.NombreCliente,
			p.MontoAporte,
			p.PorcentajeAporte,
			p.PorcentajeMeta,
			p.PorcentajeRecompensa,
			p.MontoBrutoAsignado,
			p.MontoRetencion,
			p.MontoNetoAsignado,
		); err != nil {
			return fmt.Errorf("error al guardar el detalle del reparto: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar el reparto: %w", err)
	}

	log.Printf("[MySQL] - Reparto congelado: Quiniela:%d Participantes:%d Neto:%s",
		reparto.IDQuiniela, reparto.TotalParticipantes, reparto.MontoNeto.String())
	return nil
}

func (mysql *MySQL) GetRepartoByQuiniela(idQuiniela int32) (*domain.RepartoPremio, error) {
	query := `SELECT id_reparto, id_premio, id_quiniela, monto_bruto, porcentaje_retencion,
	                 monto_neto, total_recaudado, monto_meta, porcentaje_meta_alcanzado,
	                 total_participantes, fecha_reparto
	          FROM premio_repartos WHERE id_quiniela = ?`
	row := mysql.conn.FetchRow(query, idQuiniela)

	var r domain.RepartoPremio
	var idReparto int32
	err := row.Scan(&idReparto, &r.IDPremio, &r.IDQuiniela, &r.MontoBruto, &r.PorcentajeRetencion,
		&r.MontoNeto, &r.TotalRecaudado, &r.MontoMeta, &r.PorcentajeMetaAlcanzado,
		&r.TotalParticipantes, &r.FechaReparto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("error al consultar el reparto: %w", err)
	}

	queryDetalle := `SELECT id_cliente, nombre_cliente, monto_aporte, porcentaje_aporte,
	                        porcentaje_sobre_meta, porcentaje_recompensa,
	                        monto_bruto_asignado, monto_retencion, monto_neto_asignado
	                 FROM premio_reparto_participantes
	                 WHERE id_reparto = ?
	                 ORDER BY monto_aporte DESC, id_cliente ASC`

	rows, err := mysql.conn.FetchRows(queryDetalle, idReparto)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los participantes del reparto: %w", err)
	}
	defer rows.Close()

	r.Participantes = make([]domain.ParticipanteReparto, 0, r.TotalParticipantes)
	for rows.Next() {
		var p domain.ParticipanteReparto
		if err := rows.Scan(&p.IDCliente, &p.NombreCliente, &p.MontoAporte, &p.PorcentajeAporte,
			&p.PorcentajeMeta, &p.PorcentajeRecompensa, &p.MontoBrutoAsignado,
			&p.MontoRetencion, &p.MontoNetoAsignado); err != nil {
			return nil, fmt.Errorf("error al escanear participante del reparto: %w", err)
		}
		r.Participantes = append(r.Participantes, p)
	}

	return &r, nil
}
