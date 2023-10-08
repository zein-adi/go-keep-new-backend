package auth_repos_mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"strconv"
)

var roleTableName = "roles"

func NewRoleMysqlRepository() *RoleRepository {
	db, cleanup := helpers_mysql.OpenMySqlConnection()
	return &RoleRepository{
		db:        db,
		dbCleanup: cleanup,
	}
}

type RoleRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (r *RoleRepository) Cleanup() {
	r.dbCleanup()
}
func (r *RoleRepository) Get(ctx context.Context, request auth_requests.GetRequest) []*auth_entities.Role {
	q := r.getQueryFiltered(ctx, request)
	q.Skip(request.Skip)
	q.Take(request.Take)
	q.OrderBy("nama")
	rows, cleanup := q.Get()
	defer cleanup()

	var models []*auth_entities.Role
	for rows.Next() {
		model := &auth_entities.Role{}
		models = append(models, model)

		permissionText := ""
		helpers_error.PanicIfError(rows.Scan(&model.Id, &model.Nama, &model.Deskripsi, &model.Level, &permissionText))
		helpers_error.PanicIfError(json.Unmarshal([]byte(permissionText), &model.Permissions))
	}

	return models
}

func (r *RoleRepository) Count(ctx context.Context, request auth_requests.GetRequest) (count int) {
	return r.getQueryFiltered(ctx, request).Count()
}

func (r *RoleRepository) FindById(ctx context.Context, id string) (*auth_entities.Role, error) {
	q := r.getQueryFiltered(ctx, auth_requests.NewGetRequest())
	q.Where("id", "=", id)
	q.Take(1)

	rows, cleanup := q.Get()
	defer cleanup()

	if !rows.Next() {
		return nil, helpers_error.EntryNotFoundError
	}

	model := &auth_entities.Role{}
	permissionText := ""
	helpers_error.PanicIfError(rows.Scan(&model.Id, &model.Nama, &model.Deskripsi, &model.Level, &permissionText))
	helpers_error.PanicIfError(json.Unmarshal([]byte(permissionText), &model.Permissions))
	return model, nil
}

func (r *RoleRepository) CountByNama(ctx context.Context, nama string, exceptId string) (count int) {
	q := r.getQueryFiltered(ctx, auth_requests.NewGetRequest())
	q.Where("nama", "=", nama)
	q.Where("id", "!=", exceptId)
	q.Take(1)
	return q.Count()
}

func (r *RoleRepository) Insert(ctx context.Context, role *auth_entities.Role) (*auth_entities.Role, error) {
	permissions, err := json.Marshal(role.Permissions)
	helpers_error.PanicIfError(err)

	q := helpers_mysql.NewQueryBuilder(ctx, r.db, roleTableName)
	lastId, err := q.Insert(map[string]any{
		"nama":        role.Nama,
		"deskripsi":   role.Deskripsi,
		"level":       role.Level,
		"permissions": string(permissions),
	})
	if err != nil {
		return &auth_entities.Role{}, err
	}

	model := role.Copy()
	model.Id = strconv.Itoa(lastId)
	return model, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *auth_entities.Role) (affected int, err error) {
	permissions, err := json.Marshal(role.Permissions)
	helpers_error.PanicIfError(err)

	q := helpers_mysql.NewQueryBuilder(ctx, r.db, roleTableName)
	q.Where("id", "=", role.Id)
	affected = q.Update(map[string]any{
		"nama":        role.Nama,
		"deskripsi":   role.Deskripsi,
		"level":       role.Level,
		"permissions": string(permissions),
	})
	return affected, nil
}

func (r *RoleRepository) DeleteById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, r.db, roleTableName)
	q.Where("id", "=", id)

	affected = q.Delete()
	if affected == 0 {
		return 0, helpers_error.EntryNotFoundError
	}
	return affected, nil
}

func (r *RoleRepository) getQueryFiltered(ctx context.Context, request auth_requests.GetRequest) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, r.db, roleTableName)
	q.Select("id, nama, deskripsi, level, permissions")

	if request.Search != "" {
		q.Where("nama", "LIKE", "%"+request.Search+"%")
	}

	return q
}
