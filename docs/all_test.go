package docs_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/go-accounts/docs"
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mCodec"
	"math/rand"
)

type Tester struct {
	Mango     *mango.Mango
	DocSystem *docs.DocSystem
}

func (x *Tester) Context() (ctx context.Context) {
	ctx = context.TODO()
	ctx = mango.SetMango(ctx, x.Mango)
	ctx = docs.SetJwtSecret(ctx, fnEnv.Read("JWT_SECRET"))
	ctx = docs.SetJwtIssuer(ctx, "all_test")
	ctx = docs.SetDocSystem(ctx, x.DocSystem)
	return
}

func (x *Tester) NewAccount() (res *docs.Account) {
	res = &docs.Account{
		IsActivate: true,
		Identifier: map[string]string{
			"email":    gofakeit.Email(),
			"username": gofakeit.Username(),
		},
		Property: map[string]string{
			"nickname": gofakeit.BeerName(),
			"age":      fmt.Sprintf("%d", rand.Int63n(60)),
		},
		Verifier: map[string]*docs.AccountVerifier{
			"password": {
				Key:   "salt",
				Value: "saltedPasswd",
				Mode:  docs.VerifyModeCompare,
			},
		},
		Data: make([]byte, 0),
	}

	return
}

func NewTester(truncates ...bool) (res *Tester) {
	fnPanic.On(fnEnv.ReadFromFile("../env/.env"))

	res = &Tester{}
	res.Mango = fnPanic.Get(mango.NewMango(&mango.IConn{
		Host:        fnEnv.Read("MG_HOST"),
		Username:    fnEnv.Read("MG_USERNAME"),
		Password:    fnEnv.Read("MG_PASSWORD"),
		Database:    fnEnv.Read("MG_DATABASE"),
		SetRegistry: mCodec.RegisterDecimal,
	}))

	var ctx = res.Context()
	if fnParams.Get(truncates) {
		fnPanic.On(res.Mango.DB.Drop(ctx))
	}

	fnPanic.On(res.Mango.Migrate(
		ctx,
		&docs.DocAccount{},
		&docs.DocSystem{},
		&docs.DocSession{},
		&mango.DocKv[any]{},
	))

	res.DocSystem = fnPanic.Get(docs.ReadSystem(ctx))

	return
}
