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

var userEntityName = "user"
var userTableName = "user_users"

func NewUserMysqlRepository() *UserMysqlRepository {
	db, cleanup := helpers_mysql.OpenMySqlConnection()
	return &UserMysqlRepository{
		db:        db,
		dbCleanup: cleanup,
	}
}

type UserMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (r *UserMysqlRepository) Get(ctx context.Context, request *helpers_requests.Get) []*auth_entities.User {
	q := r.getQueryFiltered(ctx, request)
	q.Skip(request.Skip)
	q.Take(request.Take)
	q.OrderBy("nama")
	return r.newEntitiesFromRows(q.Get())
}
func (r *UserMysqlRepository) Count(ctx context.Context, request *helpers_requests.Get) (count int) {
	return r.getQueryFiltered(ctx, request).Count()
}
func (r *UserMysqlRepository) FindById(ctx context.Context, id string) (*auth_entities.User, error) {
	q := r.getQueryFiltered(ctx, helpers_requests.NewGet())
	q.Where("id", "=", id)
	q.Take(1)
	models := r.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(userEntityName, "id", id)
	}
	return models[0], nil
}
func (r *UserMysqlRepository) FindByUsername(ctx context.Context, username string) (*auth_entities.User, error) {
	q := r.getQueryFiltered(ctx, helpers_requests.NewGet())
	q.Where("username", "=", username)
	q.Take(1)
	models := r.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(userEntityName, "username", username)
	}
	return models[0], nil
}
func (r *UserMysqlRepository) CountByUsername(ctx context.Context, username string, exceptId string) (count int) {
	q := r.getQueryFiltered(ctx, helpers_requests.NewGet())
	q.Where("username", "=", username)
	q.Where("id", "!=", exceptId)
	q.Take(1)
	return q.Count()
}
func (r *UserMysqlRepository) Insert(ctx context.Context, user *auth_entities.User) (*auth_entities.User, error) {
	roleIds, err := json.Marshal(user.RoleIds)
	helpers_error.PanicIfError(err)

	q := helpers_mysql.NewQueryBuilder(ctx, r.db, userTableName)
	lastId, err := q.Insert(map[string]any{
		"username": user.Username,
		"password": user.Password,
		"nama":     user.Nama,
		"role_ids": string(roleIds),
	})
	if err != nil {
		return &auth_entities.User{}, err
	}

	model := user.Copy()
	model.Id = strconv.Itoa(lastId)
	return model, nil
}
func (r *UserMysqlRepository) Update(ctx context.Context, user *auth_entities.User) (affected int, err error) {
	roleIds, err := json.Marshal(user.RoleIds)
	helpers_error.PanicIfError(err)

	q := helpers_mysql.NewQueryBuilder(ctx, r.db, userTableName)
	q.Where("id", "=", user.Id)
	affected = q.Update(map[string]any{
		"username": user.Username,
		"nama":     user.Nama,
		"role_ids": string(roleIds),
	})
	return affected, nil
}
func (r *UserMysqlRepository) UpdatePassword(ctx context.Context, userId, password string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, r.db, userTableName)
	q.Where("id", "=", userId)
	affected = q.Update(map[string]any{
		"password": password,
	})
	return affected, nil
}
func (r *UserMysqlRepository) DeleteById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, r.db, userTableName)
	q.Where("id", "=", id)

	affected = q.Delete()
	if affected == 0 {
		return 0, helpers_error.NewEntryNotFoundError(userEntityName, "id", id)
	}
	return affected, nil
}

func (r *UserMysqlRepository) Cleanup() {
	r.dbCleanup()
}
func (r *UserMysqlRepository) getQueryFiltered(ctx context.Context, request *helpers_requests.Get) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, r.db, userTableName)
	q.Select("id, username, password, nama, role_ids")

	if request.Search != "" {
		q.Where("nama", "LIKE", "%"+request.Search+"%")
	}

	return q
}
func (r *UserMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*auth_entities.User {
	defer cleanup()

	models := make([]*auth_entities.User, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &auth_entities.User{}
		models = append(models, model)

		roleIdText := ""
		helpers_error.PanicIfError(rows.Scan(
			&model.Id,
			&model.Username,
			&model.Password,
			&model.Nama,
			&roleIdText,
		))
		helpers_error.PanicIfError(json.Unmarshal([]byte(roleIdText), &model.RoleIds))
	}
	return models
}
