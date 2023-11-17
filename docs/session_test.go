package docs_test

import (
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/go-accounts/docs"
	"github.com/d3v-friends/go-pure/fnPanic"
	"testing"
)

func TestSession(test *testing.T) {
	var tr = NewTester()

	test.Run("create account, session", func(t *testing.T) {
		var ctx = tr.Context()
		var accountData = tr.NewAccount()

		var account = fnPanic.Get(docs.CreateAccount(ctx, &docs.ICreateAccount{
			Data: accountData,
		}))

		var session, token, err = docs.CreateSession(ctx, &docs.ICreateSession{
			AccountId: account.Id,
			Ip:        gofakeit.IPv4Address(),
			UserAgent: gofakeit.UserAgent(),
		})

		if err != nil {
			t.Fatal(err)
		}

	})
}
