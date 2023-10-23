package basic_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
)

func NewChangelogServices(changelogRepo basic_repo_interfaces.IChangelogRepository) *ChangelogServices {
	return &ChangelogServices{
		repo: changelogRepo,
	}
}

type ChangelogServices struct {
	repo basic_repo_interfaces.IChangelogRepository
}

func (x *ChangelogServices) Get(ctx context.Context, skip, take int) []*basic_entities.Changelog {
	return x.repo.Get(ctx, skip, take)
}
func (x *ChangelogServices) Insert(ctx context.Context, changelog *basic_entities.Changelog) (*basic_entities.Changelog, error) {
	err := validator.New().ValidateStruct(changelog)
	if err != nil {
		return nil, err
	}
	return x.repo.Insert(ctx, changelog)
}
func (x *ChangelogServices) Update(ctx context.Context, changelog *basic_entities.Changelog) (affected int, err error) {
	err = validator.New().ValidateStruct(changelog)
	if err != nil {
		return 0, err
	}
	return x.repo.Update(ctx, changelog)
}
func (x *ChangelogServices) DeleteById(ctx context.Context, id string) (affected int, err error) {
	_, err = x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	return x.repo.DeleteById(ctx, id)
}
