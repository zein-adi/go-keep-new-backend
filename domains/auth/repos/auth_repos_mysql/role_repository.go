package auth_repos_mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"strconv"
)

var roleEntityName = "role"
var roleTableName = "user_roles"

func NewRoleMysqlRepository() *RoleMysqlRepository {
	db, cleanup := helpers_mysql.OpenMySqlConnection()
	return &RoleMysqlRepository{
		db:        db,
		dbCleanup: cleanup,
	}
}

type RoleMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (r *RoleMysqlRepository) Get(ctx context.Context, request *helpers_requests.Get) []*auth_entities.Role {
	q := r.getQueryFiltered(ctx, request)
	if request.Take > 0 {
		q.Skip(request.Skip)
		q.Take(request.Take)
	}
	q.OrderBy("nama")
	return r.newEntitiesFromRows(q.Get())
}
func (r *RoleMysqlRepository) Count(ctx context.Context, request *helpers_requests.Get) (count int) {
	return r.getQueryFiltered(ctx, request).Count()
}
func (r *RoleMysqlRepository) FindById(ctx context.Context, id string) (*auth_entities.Role, error) {
	q := r.getQueryFiltered(ctx, helpers_requests.NewGet())
	q.Where("id", "=", id)
	q.Take(1)

	models := r.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(roleEntityName, "id", id)
	}
	return models[0], nil
}
func (r *RoleMysqlRepository) CountByNama(ctx context.Context, nama string, exceptId string) (count int) {
	q := r.getQueryFiltered(ctx, helpers_requests.NewGet())
	q.Where("nama", "=", nama)
	q.Where("id", "!=", exceptId)
	q.Take(1)
	return q.Count()
}
func (r *RoleMysqlRepository) Insert(ctx context.Context, role *auth_entities.Role) (*auth_entities.Role, error) {
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
func (r *RoleMysqlRepository) Update(ctx context.Context, role *auth_entities.Role) (affected int, err error) {
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
func (r *RoleMysqlRepository) DeleteById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, r.db, roleTableName)
	q.Where("id", "=", id)

	affected = q.Delete()
	if affected == 0 {
		return 0, helpers_error.NewEntryNotFoundError(roleEntityName, "id", id)
	}
	return affected, nil
}
func (r *RoleMysqlRepository) GetById(ctx context.Context, ids []string) ([]*auth_entities.Role, error) {
	q := r.getQueryFiltered(ctx, helpers_requests.NewGet())
	q.WhereIn("id", ids)
	models := r.newEntitiesFromRows(q.Get())
	expectedLen := len(ids)
	actualLen := len(models)
	if actualLen < expectedLen {
		return models, helpers_error.NewEntryCountMismatchError(expectedLen, actualLen)
	}
	return models, nil
}

func (r *RoleMysqlRepository) Cleanup() {
	r.dbCleanup()
}
func (r *RoleMysqlRepository) getQueryFiltered(ctx context.Context, request *helpers_requests.Get) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, r.db, roleTableName)
	q.Select("id, nama, deskripsi, level, permissions")

	if request.Search != "" {
		q.Where("nama", "LIKE", "%"+request.Search+"%")
	}

	return q
}
func (r *RoleMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*auth_entities.Role {
	defer cleanup()

	models := make([]*auth_entities.Role, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &auth_entities.Role{}
		models = append(models, model)

		permissionText := ""
		helpers_error.PanicIfError(rows.Scan(&model.Id, &model.Nama, &model.Deskripsi, &model.Level, &permissionText))
		helpers_error.PanicIfError(json.Unmarshal([]byte(permissionText), &model.Permissions))
	}

	return models
}
