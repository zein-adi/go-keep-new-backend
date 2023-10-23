package keep_service_interfaces

import "context"

type IBarangServices interface {
	UpdateBarangFromTransaksi(ctx context.Context) (affected int, err error)
}
