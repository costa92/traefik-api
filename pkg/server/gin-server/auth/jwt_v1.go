package auth

type GinJWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512, RS256, RS384 or RS512
	// Optional, default is HS256.
	SigningAlgorithm string
}

func NewGinJWTMiddleware(m *GinJWTMiddleware) (*GinJWTMiddleware, error) {
	if err := m.MiddlewareInit(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *GinJWTMiddleware) MiddlewareInit() error {
	return nil
}
