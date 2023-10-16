package middlewares

type AccessTokenClaims struct {
	Sub     string
	Name    string
	RoleIds []string
	Iat     int64
	Exp     int64
}

type RefreshTokenClaims struct {
	Sub      string
	Remember bool
	Iat      int64
	Exp      int64
}
