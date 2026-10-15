package domain

import "time"

type IEdicion interface {
	SaveEdicion(tipoEdicion string, nombreEdicion string, fechaLimite *time.Time) (*Edicion, error)
	DeleteEdicion(id int32) error
	UpdateEdicion(id int32, tipoEdicion string, nombreEdicion string, fechaLimite *time.Time) error
	GetAll() ([]Edicion, error)
	GetEdicionByID(id int32) (*Edicion, error)
}

type Edicion struct {
	ID            int32      `json:"id_edicion"`
	TipoEdicion   string     `json:"tipo_edicion"`   // 'MEDIA_SEMANA' o 'FIN_DE_SEMANA'
	NombreEdicion string     `json:"nombre_edicion"` // Ej: 'Jornada 10'
	FechaInicio   time.Time  `json:"fecha_inicio"`
	FechaLimite   *time.Time `json:"fecha_limite"` // cierre del concurso: ya no se aceptan aportes
}

func NewEdicion(tipoEdicion string, nombreEdicion string, fechaLimite *time.Time) *Edicion {
	return &Edicion{
		TipoEdicion:   tipoEdicion,
		NombreEdicion: nombreEdicion,
		FechaInicio:   time.Now(),
		FechaLimite:   fechaLimite,
	}
}
