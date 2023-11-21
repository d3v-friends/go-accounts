package docs_test

import (
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/go-accounts/docs"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession(test *testing.T) {
	var tr = NewTester(true)

	test.Run("create account, session", func(t *testing.T) {
		var ctx = tr.Context()
		var accountData = tr.NewAccount()

		var account = fnPanic.Get(docs.CreateAccount(ctx, &docs.ICreateAccount{
			Data: accountData,
		}))

		var iCreateSession = &docs.ICreateSession{
			AccountId: account.Id,
			Ip:        gofakeit.IPv4Address(),
			UserAgent: gofakeit.UserAgent(),
		}

		var session, token, err = docs.CreateSession(ctx, iCreateSession)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, account.Id, session.AccountId)

		var sessAccount = fnPanic.Get(docs.VerifySession(ctx, &docs.IVerifySession{
			Token:     token,
			Ip:        iCreateSession.Ip,
			UserAgent: iCreateSession.UserAgent,
		}))

		assert.Equal(t, account.Id, sessAccount.Id)
	})
}
