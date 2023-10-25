//go:build wireinject
// +build wireinject

package dependency_injection

import (
	"github.com/google/wire"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_services"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_mysql"
)

var (
	KeepPosSet = wire.NewSet(
		keep_services.NewPosServices,
		wire.Bind(new(keep_repo_interfaces.IPosRepository), new(*keep_repos_mysql.PosMysqlRepository)),
		keep_repos_mysql.NewPosMySqlRepository,
	)
	KeepKantongSet = wire.NewSet(
		keep_services.NewKantongServices,
		wire.Bind(new(keep_repo_interfaces.IKantongRepository), new(*keep_repos_mysql.KantongMysqlRepository)),
		keep_repos_mysql.NewKantongMysqlRepository,
	)
	KeepKantongHistorySet = wire.NewSet(
		keep_services.NewKantongHistoryServices,
		wire.Bind(new(keep_repo_interfaces.IKantongHistoryRepository), new(*keep_repos_mysql.KantongHistoryMysqlRepository)),
		keep_repos_mysql.NewKantongHistoryMysqlRepository,
	)
	KeepTransaksiSet = wire.NewSet(
		keep_services.NewTransaksiServices,
		wire.Bind(new(keep_repo_interfaces.ITransaksiRepository), new(*keep_repos_mysql.TransaksiMysqlRepository)),
		keep_repos_mysql.NewTransaksiMySqlRepository,
	)
	KeepLokasiSet = wire.NewSet(
		keep_services.NewLokasiServices,
		wire.Bind(new(keep_repo_interfaces.ILokasiRepository), new(*keep_repos_mysql.LokasiMysqlRepository)),
		keep_repos_mysql.NewLokasiMySqlRepository,
	)
	KeepBarangSet = wire.NewSet(
		keep_services.NewBarangServices,
		wire.Bind(new(keep_repo_interfaces.IBarangRepository), new(*keep_repos_mysql.BarangMysqlRepository)),
		keep_repos_mysql.NewBarangMySqlRepository,
	)
)

func InitKeepPosServices() *keep_services.PosServices {
	wire.Build(KeepPosSet, KeepTransaksiSet)
	return nil
}
func InitKeepKantongServices() *keep_services.KantongServices {
	wire.Build(KeepKantongSet, KeepPosSet)
	return nil
}
func InitKeepKantongHistoryServices() *keep_services.KantongHistoryServices {
	wire.Build(KeepKantongHistorySet, KeepKantongSet)
	return nil
}
func InitKeepTransaksiServices() *keep_services.TransaksiServices {
	wire.Build(KeepTransaksiSet, KeepPosSet, KeepKantongSet)
	return nil
}
func InitKeepLokasiServices() *keep_services.LokasiServices {
	wire.Build(KeepLokasiSet, KeepTransaksiSet)
	return nil
}
func InitKeepBarangServices() *keep_services.BarangServices {
	wire.Build(KeepBarangSet, KeepTransaksiSet)
	return nil
}
