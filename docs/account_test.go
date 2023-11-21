package docs_test

import (
	"fmt"
	"github.com/d3v-friends/go-accounts/docs"
	"github.com/d3v-friends/go-pure/fnMatch"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mType"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestAccount(test *testing.T) {
	var tr = NewTester(true)

	test.Run("reindexing", func(t *testing.T) {
		var ctx = tr.Context()
		var allIdx = []string{
			"_id_",
			"data.identifer.email_1",
			"data.identifier.username_1",
			"data.property.nickname_1",
			"data.property.age_1",
		}

		fnPanic.On(docs.ReindexingAccount(ctx, &docs.IReindexingAccount{
			IdentifierKeys: []string{"email", "username"},
			PropertyKeys:   []string{"nickname", "age"},
		}))

		var account = &docs.DocAccount{}

		var err error
		var col = mango.GetColP(ctx, account.GetColNm())
		var cur *mongo.Cursor
		if cur, err = col.Indexes().List(ctx); err != nil {
			t.Fatal(err)
		}

		var ls = make(mType.IndexModels, 0)
		if err = cur.All(ctx, &ls); err != nil {
			t.Fatal(err)
		}

		if len(ls) != 5 {
			t.Fatal(fmt.Errorf(
				"invalid index count: expected=%d, count=%d",
				5,
				len(ls),
			))
		}

		for _, idName := range allIdx {
			_, err = fnMatch.Get(ls, func(v *mType.IndexModel) bool {
				return v.Name == idName
			})

			if err != nil {
				t.Fatal(fmt.Errorf("not found index: idx_nm=%s", idName))
			}
		}

	})

	test.Run("create account", func(t *testing.T) {
		var ctx = tr.Context()
		var data = tr.NewAccount()
		var account = fnPanic.Get(docs.CreateAccount(ctx, &docs.ICreateAccount{
			Data: data,
		}))

		assert.Equal(t, data.Identifier["email"], account.Data.Identifier["email"])
		assert.Equal(t, data.Identifier["username"], account.Data.Identifier["username"])
	})
}
