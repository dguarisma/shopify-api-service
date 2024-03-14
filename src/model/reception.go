package model

import (
	"time"
)

type ReceptionArt struct { // articulos recibidos
	// no puede cambiarse
	CustomModel
	ArticleID uint

	Missing uint  // faltantes
	Refund  uint  // devoluciones
	Reason  uint8 // [1, 2, 3]
	Count   uint
	Batch   string // lote
	Date    time.Time
}
