package keep_service_interfaces

import "context"

type ILokasiServices interface {
	UpdateLokasiFromTransaksi(ctx context.Context) (affected int, err error)
}
