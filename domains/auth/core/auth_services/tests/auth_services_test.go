package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_memory"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	r := AuthServicesTest{}

	t.Run("LoginSuccess", func(t *testing.T) {
		r.setup()
		r.populateRole()
		r.populateUser()

		username := "zeinadimukadar"
		password := "aA123456789"

		accessToken, refreshToken, err := r.services.Login(context.Background(), username, password, true)
		assert.Nil(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)

		user, err := r.userRepo.FindByUsername(context.Background(), username)
		assert.Nil(t, err)

		claims, err := middlewares.GetJwtClaims(accessToken)
		roleIds := claims.RoleIds

		assert.Nil(t, err)
		assert.Nil(t, err)
		assert.Equal(t, username, claims.Sub)
		assert.Equal(t, user.Nama, claims.Name)
		assert.Equal(t, user.RoleIds, roleIds)

		refreshClaims, err := middlewares.GetRefreshJwtClaims(refreshToken)
		assert.Nil(t, err)
		assert.Nil(t, err)
		assert.Equal(t, username, refreshClaims.Sub)
	})
	t.Run("LoginSuccessRememberMeTrue", func(t *testing.T) {
		r.setup()
		r.populateRole()
		r.populateUser()

		username := "zeinadimukadar"
		password := "aA123456789"

		_, refreshToken, _ := r.services.Login(context.Background(), username, password, true)

		refreshClaims, _ := middlewares.GetRefreshJwtClaims(refreshToken)
		expired := time.Unix(refreshClaims.Exp, 0)
		assert.Equal(t, true, refreshClaims.Remember)

		// Not Expired Yet
		nowPlus6Days23Hour := time.Now().Add(time.Hour*24*7 - time.Hour*1)
		assert.False(t, expired.Before(nowPlus6Days23Hour))
		// Expired
		nowPlus7Days1Minute := time.Now().Add(time.Hour*24*7 + time.Minute*1)
		assert.True(t, expired.Before(nowPlus7Days1Minute))

		willBeExpiredAtLeast7Days := time.Now().Add(time.Hour*24*7 - 1)
		assert.True(t, expired.Before(willBeExpiredAtLeast7Days))
	})
	t.Run("LoginSuccessRememberMeFalse", func(t *testing.T) {
		r.setup()
		r.populateRole()
		r.populateUser()

		username := "zeinadimukadar"
		password := "aA123456789"

		_, refreshToken, _ := r.services.Login(context.Background(), username, password, false)

		refreshClaims, _ := middlewares.GetRefreshJwtClaims(refreshToken)
		expired := time.Unix(refreshClaims.Exp, 0)
		assert.Equal(t, false, refreshClaims.Remember)

		// Not Expired Yet
		nowPlus23Hour59Minute := time.Now().Add(time.Hour*24 - time.Minute*1)
		assert.False(t, expired.Before(nowPlus23Hour59Minute))
		// Expired
		nowPlus24Hour1Minute := time.Now().Add(time.Hour*24 + time.Minute*1)
		assert.True(t, expired.Before(nowPlus24Hour1Minute))
	})
	t.Run("LoginFailedCauseValidation", func(t *testing.T) {
		r.setup()

		tests := []map[string]string{
			{
				"u": "",
				"p": "",
				"e": "username.required",
			},
			{
				"u": "",
				"p": "",
				"e": "password.required",
			},
		}
		for _, test := range tests {
			_, _, err := r.services.Login(context.Background(), test["u"], test["p"], true)
			assert.ErrorContains(t, err, test["e"])
		}
	})
	t.Run("RefreshTokenFailedCauseBlacklisted", func(t *testing.T) {
		r.setup()
		r.populateRole()
		user := r.populateUser()[0]

		iat := time.Now()
		exp := time.Now().Add(time.Hour)
		refreshToken := middlewares.IssueNewRefreshToken(user.Username, false, iat, exp)
		_ = r.services.Logout(context.Background(), refreshToken)

		newAcc, newRefresh, err := r.services.Refresh(context.Background(), refreshToken)
		assert.NotNil(t, err)
		assert.Empty(t, newAcc)
		assert.Empty(t, newRefresh)
	})
	t.Run("LogoutSuccess", func(t *testing.T) {
		r.setup()

		refreshToken := middlewares.IssueNewRefreshToken("zeinadimukadar", false, time.Now(), time.Now().Add(time.Hour))
		err := r.services.Logout(context.Background(), refreshToken)
		assert.Nil(t, err)

		entryNotFoundError := r.repo.FindBlacklistByToken(context.Background(), refreshToken)
		assert.Nil(t, entryNotFoundError)
	})
	t.Run("ProfileSuccess", func(t *testing.T) {
		r.setup()
		r.populateRole()
		users := r.populateUser()
		// Bambang Santoso, Role: Guru dan Staf
		user := users[3]

		refreshToken := middlewares.IssueNewAccessToken(user.Username, user.Nama, user.RoleIds, time.Now(), time.Now().Add(time.Hour))
		profile, err := r.services.Profile(context.Background(), refreshToken)
		assert.Nil(t, err)
		assert.Equal(t, user.Username, profile.Username)
		assert.Equal(t, user.Nama, profile.Nama)

		assert.Contains(t, profile.Roles, "Guru")
		assert.Contains(t, profile.Roles, "Staf")
		assert.Contains(t, profile.Permissions, "master.siswa.get")
		assert.Contains(t, profile.Permissions, "kbm.kelas.get")
		assert.Contains(t, profile.Permissions, "master.staf.get")
	})
	t.Run("RefreshTokenSuccessSame", func(t *testing.T) {
		r.setup()
		r.populateRole()
		users := r.populateUser()
		user := users[0]

		// Issued 20 hour ago
		iat := time.Now().Add(-time.Hour * 20)
		// Will Be Expired in 4 hours
		exp := time.Now().Add(time.Hour * 4)
		refreshToken := middlewares.IssueNewRefreshToken(user.Username, false, iat, exp)
		newAccToken, newRefToken, err := r.services.Refresh(context.Background(), refreshToken)
		assert.Nil(t, err)
		assert.NotEmpty(t, newAccToken)
		assert.Equal(t, refreshToken, newRefToken)

		// Issued 6 days ago
		iat = time.Now().Add(-time.Hour * 24 * 6)
		// Will Be Expired in 25 hour
		exp = time.Now().Add(time.Hour * 25)
		refreshToken = middlewares.IssueNewRefreshToken(user.Username, true, iat, exp)
		newAccToken, newRefToken, err = r.services.Refresh(context.Background(), refreshToken)
		assert.Nil(t, err)
		assert.NotEmpty(t, newAccToken)
		assert.Equal(t, refreshToken, newRefToken)
	})
	t.Run("RefreshTokenSuccessChanged", func(t *testing.T) {
		r.setup()
		r.populateRole()
		users := r.populateUser()
		user := users[0]

		// Issued 22 hour ago
		iat := time.Now().Add(-time.Hour * 22)
		// Will Be Expired in 2 hours
		exp := time.Now().Add(time.Hour * 2)
		refreshToken := middlewares.IssueNewRefreshToken(user.Username, false, iat, exp)
		newAccToken, newRefToken, err := r.services.Refresh(context.Background(), refreshToken)
		assert.Nil(t, err)
		assert.NotEmpty(t, newAccToken)
		assert.NotEqual(t, refreshToken, newRefToken)

		// Issued 6 days ago
		iat = time.Now().Add(-time.Hour * 24 * 6)
		// Will Be Expired in 23 hour
		exp = time.Now().Add(time.Hour * 23)
		refreshToken = middlewares.IssueNewRefreshToken(user.Username, true, iat, exp)
		newAccToken, newRefToken, err = r.services.Refresh(context.Background(), refreshToken)
		assert.Nil(t, err)
		assert.NotEmpty(t, newAccToken)
		assert.NotEqual(t, refreshToken, newRefToken)
	})
	t.Run("LogoutProfileRefreshTokenFailedTokenExpired", func(t *testing.T) {
		r.setup()
		r.populateRole()
		users := r.populateUser()
		user := users[0]

		// Issued 24 hour ago
		iat := time.Now().Add(-time.Hour * 24)
		// Expired 1 minute
		exp := time.Now().Add(-time.Minute)
		refreshToken := middlewares.IssueNewRefreshToken(user.Username, false, iat, exp)

		err := r.services.Logout(context.Background(), refreshToken)
		assert.NotNil(t, err)

		_, err = r.services.Profile(context.Background(), refreshToken)
		assert.NotNil(t, err)

		newAccToken, newRefToken, err := r.services.Refresh(context.Background(), refreshToken)
		assert.NotNil(t, err)
		assert.Empty(t, newAccToken)
		assert.Empty(t, newRefToken)
	})
	t.Run("ProfileEmptyWhenRolesNotFound", func(t *testing.T) {
		r.setup()
		r.populateRole()
		users := r.populateUser()
		// Bambang Santoso, Role: Guru dan Staf
		user := users[3]
		user.RoleIds = []string{"10000"}

		refreshToken := middlewares.IssueNewAccessToken(user.Username, user.Nama, user.RoleIds, time.Now(), time.Now().Add(time.Hour))
		profile, err := r.services.Profile(context.Background(), refreshToken)
		assert.Nil(t, err)
		assert.Equal(t, user.Username, profile.Username)
		assert.Equal(t, user.Nama, profile.Nama)
		assert.Empty(t, profile.Roles)
		assert.Empty(t, profile.Permissions)
	})
}

type AuthServicesTest struct {
	repo     auth_repo_interfaces.IAuthRepository
	services auth_service_interfaces.IAuthServices
	userRepo auth_repo_interfaces.IUserRepository
	roleRepo auth_repo_interfaces.IRoleRepository
}

func (r *AuthServicesTest) setup() {
	r.setMemoryRepository()
	r.services = auth_services.NewAuthServices(r.repo, r.userRepo, r.roleRepo)
}
func (r *AuthServicesTest) setMemoryRepository() {
	r.repo = auth_repos_memory.NewAuthMemoryRepository()
	//r.repo = auth_repos_redis.NewAuthRedisRepository()
	r.userRepo = auth_repos_memory.NewUserMemoryRepository()
	r.roleRepo = auth_repos_memory.NewRoleMemoryRepository()
}
func (r *AuthServicesTest) populateRole() []*auth_entities.Role {
	input := []*auth_entities.Role{
		{
			Id:    "1",
			Nama:  "Developer",
			Level: 1,
			Permissions: []string{
				"user.user.get",
				"user.user.insert",
				"user.user.update",
				"user.user.delete",
			},
		},
		{
			Id:    "2",
			Nama:  "Guru",
			Level: 100,
			Permissions: []string{
				"master.siswa.get",
				"kbm.kelas.get",
			},
		},
		{
			Id:    "3",
			Nama:  "Staf",
			Level: 100,
			Permissions: []string{
				"master.staf.get",
				"master.siswa.get",
			},
		},
	}
	var models []*auth_entities.Role
	for _, in := range input {
		model, _ := r.roleRepo.Insert(context.Background(), in)
		models = append(models, model)
	}
	return models
}
func (r *AuthServicesTest) populateUser() []*auth_entities.User {
	input := []*auth_entities.User{
		{
			Id:       "1",
			Username: "zeinadimukadar",
			Password: auth_services.HashPassword("aA123456789"),
			Nama:     "Zein Adi",
			RoleIds:  []string{"1"},
		},
		{
			Id:       "2",
			Username: "witadwinor",
			Password: auth_services.HashPassword("aA1234567890"),
			Nama:     "Wita Dwi Nor",
			RoleIds:  []string{"2"},
		},
		{
			Id:       "3",
			Username: "agungsetya",
			Password: auth_services.HashPassword("aA1234567890"),
			Nama:     "Agung Setya Arifin",
			RoleIds:  []string{"3"},
		},
		{
			Id:       "4",
			Username: "bambangsantoso",
			Password: auth_services.HashPassword("aA1234567890"),
			Nama:     "Bambang Santoso",
			RoleIds:  []string{"2", "3"},
		},
	}
	var models []*auth_entities.User
	for _, in := range input {
		model, _ := r.userRepo.Insert(context.Background(), in)
		models = append(models, model)
	}
	return models
}
