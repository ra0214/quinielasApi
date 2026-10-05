package domain

import "time"

type IEdicion interface {
	SaveEdicion(tipoEdicion string, nombreEdicion string) (*Edicion, error)
	DeleteEdicion(id int32) error
	UpdateEdicion(id int32, tipoEdicion string, nombreEdicion string) error
	GetAll() ([]Edicion, error)
	GetEdicionByID(id int32) (*Edicion, error)
}

type Edicion struct {
	ID            int32     `json:"id_edicion"`
	TipoEdicion   string    `json:"tipo_edicion"`   // 'MEDIA_SEMANA' o 'FIN_DE_SEMANA'
	NombreEdicion string    `json:"nombre_edicion"` // Ej: 'Jornada 10'
	FechaInicio   time.Time `json:"fecha_inicio"`
}

func NewEdicion(tipoEdicion string, nombreEdicion string) *Edicion {
	return &Edicion{
		TipoEdicion:   tipoEdicion,
		NombreEdicion: nombreEdicion,
		FechaInicio:   time.Now(),
	}
}