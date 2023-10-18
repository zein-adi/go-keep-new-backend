//go:build wireinject
// +build wireinject

package dependency_injection

import (
	"github.com/google/wire"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_services"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/handlers/keep_handlers_restful"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_mysql"
)

var (
	KeepPosSet = wire.NewSet(
		keep_handlers_restful.NewPosRestfulHandler,
		wire.Bind(new(keep_service_interfaces.IPosServices), new(*keep_services.PosServices)),
		keep_services.NewPosServices,
		wire.Bind(new(keep_repo_interfaces.IPosRepository), new(*keep_repos_mysql.PosMysqlRepository)),
		keep_repos_mysql.NewPosMySqlRepository,
	)
	KeepKantongSet = wire.NewSet(
		keep_handlers_restful.NewKantongRestfulHandler,
		wire.Bind(new(keep_service_interfaces.IKantongServices), new(*keep_services.KantongServices)),
		keep_services.NewKantongServices,
		wire.Bind(new(keep_repo_interfaces.IKantongRepository), new(*keep_repos_mysql.KantongMysqlRepository)),
		keep_repos_mysql.NewKantongMysqlRepository,
	)
	KeepKantongHistorySet = wire.NewSet(
		keep_handlers_restful.NewKantongHistoryRestfulHandler,
		wire.Bind(new(keep_service_interfaces.IKantongHistoryServices), new(*keep_services.KantongHistoryServices)),
		keep_services.NewKantongHistoryServices,
		wire.Bind(new(keep_repo_interfaces.IKantongHistoryRepository), new(*keep_repos_mysql.KantongHistoryMysqlRepository)),
		keep_repos_mysql.NewKantongHistoryMysqlRepository,
	)
)

func InitKeepPosRestful() *keep_handlers_restful.PosRestfulHandler {
	wire.Build(KeepPosSet)
	return nil
}
func InitKeepKantongRestful() *keep_handlers_restful.KantongRestfulHandler {
	wire.Build(KeepKantongSet)
	return nil
}
func InitKeepKantongHistoryRestful() *keep_handlers_restful.KantongHistoryRestfulHandler {
	wire.Build(KeepKantongHistorySet, KeepKantongSet)
	return nil
}
