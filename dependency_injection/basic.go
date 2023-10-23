//go:build wireinject
// +build wireinject

package dependency_injection

import (
	"github.com/google/wire"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_services"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/repos/basic_repos_file"
)

var (
	BasicChangelogSet = wire.NewSet(
		basic_services.NewChangelogServices,
		wire.Bind(new(basic_repo_interfaces.IChangelogRepository), new(*basic_repos_file.ChangelogFileRepository)),
		basic_repos_file.NewChangelogFileRepository,
	)
)

func InitBasicChangelogServices() *basic_services.ChangelogServices {
	wire.Build(BasicChangelogSet)
	return nil
}
