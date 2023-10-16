package auth_services

import (
	"context"
	"errors"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_responses"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	refreshTokenBlacklistedError = errors.New("refresh token is blacklisted")
)

func NewAuthServices(authRepo auth_repo_interfaces.IAuthRepository, userRepo auth_repo_interfaces.IUserRepository, roleRepo auth_repo_interfaces.IRoleRepository) *AuthServices {
	return &AuthServices{
		authRepo: authRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

type AuthServices struct {
	authRepo auth_repo_interfaces.IAuthRepository
	userRepo auth_repo_interfaces.IUserRepository
	roleRepo auth_repo_interfaces.IRoleRepository
}

func (a *AuthServices) Login(ctx context.Context, username, rawPassword string, rememberMe bool) (accessToken, refreshToken string, err error) {
	v := validator.New()
	data := map[string]interface{}{
		"username": username,
		"password": rawPassword,
	}
	rules := map[string]interface{}{
		"username": "required",
		"password": "required",
	}
	err = v.ValidateMap(data, rules)
	if err != nil {
		return "", "", err
	}

	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", helpers_error.NewValidationErrors("credentials", "invalid", "")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPassword))
	if err != nil {
		return "", "", helpers_error.NewValidationErrors("credentials", "invalid", "")
	}

	iat := time.Now()
	exp := time.Now().Add(time.Minute * 10)
	accessToken = middlewares.IssueNewAccessToken(user.Username, user.Nama, user.RoleIds, iat, exp)

	if rememberMe {
		exp = time.Now().Add(time.Hour * 24 * 7)
	} else {
		exp = time.Now().Add(time.Hour * 24 * 1)
	}
	refreshToken = middlewares.IssueNewRefreshToken(user.Username, rememberMe, iat, exp)

	return accessToken, refreshToken, nil
}

func (a *AuthServices) Refresh(ctx context.Context, refreshToken string) (accessToken, updatedRefreshToken string, err error) {
	claims, err := middlewares.GetRefreshJwtClaims(refreshToken)
	if err != nil {
		return "", "", err
	}
	isEntryFound := a.authRepo.FindBlacklistByToken(ctx, refreshToken) == nil
	if isEntryFound {
		return "", "", refreshTokenBlacklistedError
	}

	username, err := claims.GetSubject()
	if err != nil {
		return "", "", err
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return "", "", err
	}
	rememberMe := claims["remember"].(bool)

	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	newIat := time.Now()
	newExp := time.Now().Add(time.Minute * 15)
	accessToken = middlewares.IssueNewAccessToken(user.Username, user.Nama, user.RoleIds, newIat, newExp)

	var futureExpired time.Time
	if rememberMe {
		newExp = time.Now().Add(time.Hour * 24 * 7)
		futureExpired = newIat.Add(time.Hour * 24)
	} else {
		newExp = time.Now().Add(time.Hour * 24 * 1)
		futureExpired = newIat.Add(time.Hour * 3)
	}

	shouldIssueNewRefreshToken := exp.Before(futureExpired)
	if shouldIssueNewRefreshToken {
		refreshToken = middlewares.IssueNewRefreshToken(user.Username, rememberMe, newIat, newExp)
	}
	return accessToken, refreshToken, nil
}

func (a *AuthServices) Logout(ctx context.Context, refreshToken string) error {
	_, err := middlewares.GetRefreshJwtClaims(refreshToken)
	if err != nil {
		return err
	}
	return a.authRepo.InsertBlackList(ctx, refreshToken)
}

func (a *AuthServices) Profile(ctx context.Context, accessToken string) (*auth_responses.ProfileResponse, error) {
	response := &auth_responses.ProfileResponse{}
	claims, err := middlewares.GetJwtClaims(accessToken)
	if err != nil {
		return response, err
	}

	username, err := claims.GetSubject()
	if err != nil {
		return response, err
	}

	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return response, err
	}
	roles, _ := a.roleRepo.GetById(ctx, user.RoleIds)
	rolesKeyId := helpers.KeyBy(roles, func(d *auth_entities.Role) string {
		return d.Id
	})

	var roleNames []string
	var permissions []string
	for _, roleId := range user.RoleIds {
		role := rolesKeyId[roleId]
		roleNames = append(roleNames, role.Nama)
		permissions = append(permissions, role.Permissions...)
	}
	permissions = helpers.Unique(permissions)

	response = &auth_responses.ProfileResponse{
		Username:    user.Username,
		Nama:        user.Nama,
		Roles:       roleNames,
		Permissions: permissions,
	}
	return response, nil
}
