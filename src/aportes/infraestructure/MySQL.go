package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"quinielas/src/aportes/domain"
	"quinielas/src/config"
	"quinielas/src/shared/money"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IAporte = (*MySQL)(nil)

func NewMySQL() domain.IAporte {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveOrUpdateAporte(idCliente int32, idQuiniela int32, montoSumar money.Money, montoMetaQuiniela money.Money) (*domain.Aporte, error) {
	// 1. Verificar si el cliente ya cuenta con un aporte previo en esta quiniela
	querySelect := "SELECT id_aporte, monto_total_acumulado FROM aportes WHERE id_cliente = ? AND id_quiniela = ?"
	row := mysql.conn.FetchRow(querySelect, idCliente, idQuiniela)

	var idAporte int32
	var montoAcumulado money.Money
	err := row.Scan(&idAporte, &montoAcumulado)

	var nuevoTotal money.Money
	if err != nil {
		if err == sql.ErrNoRows {
			// No existe registro: Insertar nuevo aporte
			nuevoTotal = montoSumar
			porcentaje := money.PercentOf(nuevoTotal, montoMetaQuiniela)

			queryInsert := "INSERT INTO aportes (id_cliente, id_quiniela, monto_total_acumulado, porcentaje_participacion) VALUES (?, ?, ?, ?)"
			res, errIns := mysql.conn.ExecutePreparedQuery(queryInsert, idCliente, idQuiniela, nuevoTotal, porcentaje)
			if errIns != nil {
				return nil, fmt.Errorf("error al insertar aporte: %v", errIns)
			}

			lastID, _ := res.LastInsertId()
			aporte := domain.NewAporte(idCliente, idQuiniela, nuevoTotal, porcentaje)
			aporte.IDAporte = int32(lastID)
			return aporte, nil
		}
		return nil, fmt.Errorf("error al verificar aporte existente: %v", err)
	}

	// Si ya existía: Sumar al acumulado y recalcular porcentaje
	nuevoTotal = montoAcumulado.Add(montoSumar)
	porcentaje := money.PercentOf(nuevoTotal, montoMetaQuiniela)

	queryUpdate := "UPDATE aportes SET monto_total_acumulado = ?, porcentaje_participacion = ? WHERE id_aporte = ?"
	_, errUpd := mysql.conn.ExecutePreparedQuery(queryUpdate, nuevoTotal, porcentaje, idAporte)
	if errUpd != nil {
		return nil, fmt.Errorf("error al actualizar aporte: %v", errUpd)
	}

	aporte := domain.NewAporte(idCliente, idQuiniela, nuevoTotal, porcentaje)
	aporte.IDAporte = idAporte
	return aporte, nil
}

func (mysql *MySQL) UpdateAporte(idAporte int32, montoNuevo money.Money, montoMetaQuiniela money.Money) (*domain.Aporte, error) {
	aporte, err := mysql.GetAporteByID(idAporte)
	if err != nil {
		return nil, err
	}

	porcentaje := money.PercentOf(montoNuevo, montoMetaQuiniela)
	query := "UPDATE aportes SET monto_total_acumulado = ?, porcentaje_participacion = ? WHERE id_aporte = ?"
	if _, err := mysql.conn.ExecutePreparedQuery(query, montoNuevo, porcentaje, idAporte); err != nil {
		return nil, fmt.Errorf("error al actualizar aporte: %v", err)
	}

	resultado := domain.NewAporte(aporte.IDCliente, aporte.IDQuiniela, montoNuevo, porcentaje)
	resultado.IDAporte = idAporte
	return resultado, nil
}

func (mysql *MySQL) DeleteAporte(idAporte int32) error {
	query := "DELETE FROM aportes WHERE id_aporte = ?"
	_, err := mysql.conn.ExecutePreparedQuery(query, idAporte)
	if err != nil {
		return fmt.Errorf("error al eliminar aporte: %v", err)
	}
	return nil
}

func (mysql *MySQL) GetAporteByID(idAporte int32) (*domain.Aporte, error) {
	query := "SELECT id_aporte, id_cliente, id_quiniela, monto_total_acumulado, porcentaje_participacion FROM aportes WHERE id_aporte = ?"
	row := mysql.conn.FetchRow(query, idAporte)

	var a domain.Aporte
	err := row.Scan(&a.IDAporte, &a.IDCliente, &a.IDQuiniela, &a.MontoTotalAcumulado, &a.PorcentajeParticipacion)
	if err != nil {
		return nil, fmt.Errorf("aporte no encontrado: %v", err)
	}

	return &a, nil
}

func (mysql *MySQL) GetAportesByQuinielaID(idQuiniela int32) ([]domain.AporteDetalle, error) {
	query := `SELECT a.id_aporte, a.id_cliente, c.nombre, a.id_quiniela, a.monto_total_acumulado, a.porcentaje_participacion 
	          FROM aportes a 
	          JOIN clientes c ON a.id_cliente = c.id_cliente 
	          WHERE a.id_quiniela = ? 
	          ORDER BY a.monto_total_acumulado DESC`

	rows, err := mysql.conn.FetchRows(query, idQuiniela)
	if err != nil {
		return nil, fmt.Errorf("error al obtener aportes de la quiniela: %v", err)
	}
	defer rows.Close()

	var lista []domain.AporteDetalle
	for rows.Next() {
		var ad domain.AporteDetalle
		if err := rows.Scan(&ad.IDAporte, &ad.IDCliente, &ad.NombreCliente, &ad.IDQuiniela, &ad.MontoTotalAcumulado, &ad.PorcentajeParticipacion); err != nil {
			return nil, fmt.Errorf("error al escanear fila de aporte: %v", err)
		}
		lista = append(lista, ad)
	}

	return lista, nil
}

func (mysql *MySQL) GetAportesByClienteID(idCliente int32) ([]domain.Aporte, error) {
	query := "SELECT id_aporte, id_cliente, id_quiniela, monto_total_acumulado, porcentaje_participacion FROM aportes WHERE id_cliente = ?"
	rows, err := mysql.conn.FetchRows(query, idCliente)
	if err != nil {
		return nil, fmt.Errorf("error al obtener aportes del cliente: %v", err)
	}
	defer rows.Close()

	var lista []domain.Aporte
	for rows.Next() {
		var a domain.Aporte
		if err := rows.Scan(&a.IDAporte, &a.IDCliente, &a.IDQuiniela, &a.MontoTotalAcumulado, &a.PorcentajeParticipacion); err != nil {
			return nil, fmt.Errorf("error al escanear fila de aporte: %v", err)
		}
		lista = append(lista, a)
	}

	return lista, nil
}