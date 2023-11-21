package docs

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-pure/fnCtx"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mMigrate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
)

const (
	docSystem       = "systems"
	keySession      = "session"
	keyAccountIndex = "accountIndex"
)

type (
	DocSystem struct{}

	KvSession struct {
		JwtSecret      string        `bson:"jwtSecret"`
		Issuer         string        `bson:"issuer"`
		ExpireAt       time.Duration `bson:"expireAt"`
		CheckIp        bool          `bson:"checkIp"`
		CheckUserAgent bool          `bson:"checkUserAgent"`
	}

	KvAccountIndex struct {
		Identifier []string
		Property   []string
	}
)

func (x *DocSystem) GetID() primitive.ObjectID {
	return primitive.NilObjectID
}

func (x *DocSystem) GetColNm() string {
	return docSystem
}

func (x *DocSystem) GetMigrateList() mMigrate.FnMigrateList {
	return mgSystem
}

var mgSystem = mMigrate.FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "create system key_values"

		if _, err = mango.SetKv[KvSession](ctx, keySession, &KvSession{
			JwtSecret:      fmt.Sprintf("%12d", rand.Int63()),
			Issuer:         "go-accounts",
			ExpireAt:       -1,
			CheckIp:        false,
			CheckUserAgent: true,
		}); err != nil {
			return
		}

		if _, err = mango.SetKv[KvAccountIndex](ctx, keyAccountIndex, &KvAccountIndex{
			Identifier: []string{},
			Property:   []string{},
		}); err != nil {
			return
		}

		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */

func GetKvSession(ctx context.Context) (res *KvSession, err error) {
	return mango.GetKv[KvSession](ctx, keySession)
}

func SetKvSession(ctx context.Context, v *KvSession) (err error) {
	_, err = mango.SetKv[KvSession](ctx, keySession, v)
	return
}

/* ------------------------------------------------------------------------------------------------------------ */

const ctxKvSession = "CTX_KV_SESSION"

var GetCtxKvSession = fnCtx.GetFnP[*KvSession](ctxKvSession)
var SetCtxKvSession = fnCtx.SetFn[*KvSession](ctxKvSession)

/* ------------------------------------------------------------------------------------------------------------ */
