package fnJWT

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type V1[DATA IfJwtData] struct {
	secret []byte
	issuer string
	expire time.Duration
}

type IfJwtData interface {
	GetSessionID() string
	GetAudience() []string
	GetSubject() string
}

func NewV1[DATA IfJwtData](
	secret, issuer string,
	expires ...time.Duration,
) *V1[DATA] {
	var expire = fnParams.Get(expires)
	if expire == 0 {
		expire = -1
	}

	return &V1[DATA]{
		secret: []byte(secret),
		issuer: issuer,
		expire: expire,
	}
}

func (x V1[DATA]) Encode(data DATA) (res string, err error) {
	var now = time.Now()
	var nowNumericDate = &jwt.NumericDate{
		Time: now,
	}

	var claims = &jwt.RegisteredClaims{
		Issuer:    x.issuer,
		Subject:   data.GetSubject(),
		Audience:  data.GetAudience(),
		NotBefore: nowNumericDate,
		IssuedAt:  nowNumericDate,
		ID:        data.GetSessionID(),
	}

	if 0 < x.expire {
		claims.ExpiresAt = &jwt.NumericDate{
			Time: now.Add(x.expire),
		}
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	if res, err = token.SignedString(x.secret); err != nil {
		return
	}

	return
}

func (x V1[DATA]) Decode(str string) (claims *jwt.RegisteredClaims, err error) {
	claims = new(jwt.RegisteredClaims)
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(str, claims, func(token *jwt.Token) (interface{}, error) {
		return x.secret, nil
	}); err != nil {
		return
	}

	if !token.Valid {
		claims = nil
		err = fmt.Errorf("invalid jwt_str: str=%s", str)
		return
	}

	return
}
