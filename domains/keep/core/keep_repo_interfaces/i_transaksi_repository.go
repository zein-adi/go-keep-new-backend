package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
)

type ITransaksiRepository interface {
	Insert(ctx context.Context, transaksi *keep_entities.Transaksi) (*keep_entities.Transaksi, error)
}
