//go:build wireinject
// +build wireinject

package dependency_injection

import (
	"github.com/google/wire"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_services"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/handlers/keep_handlers_restful"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/handlers/keep_handlers_restful/keep_handlers_restful_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_mysql"
)

var (
	KeepPosSet = wire.NewSet(
		keep_repos_mysql.NewPosMySqlRepository,
		wire.Bind(new(keep_repo_interfaces.IPosRepository), new(*keep_repos_mysql.PosMysqlRepository)),
		keep_services.NewPosServices,
		wire.Bind(new(keep_service_interfaces.IPosServices), new(*keep_services.PosServices)),
		keep_handlers_restful.NewPosRestfulHandler,
		wire.Bind(new(keep_handlers_restful_interfaces.IPosRestfulHandler), new(*keep_handlers_restful.PosRestfulHandler)),
	)
)

func InitKeepPosRestful() *keep_handlers_restful.PosRestfulHandler {
	wire.Build(KeepPosSet)
	return nil
}
