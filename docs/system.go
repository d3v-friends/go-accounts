package docs

import (
	"context"
	"github.com/d3v-friends/go-pure/fnCtx"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mMigrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	docSystem = "systems"
)

var (
	docSystemId = primitive.NilObjectID
)

type (
	DocSystem mango.MDoc[System]

	System struct {
		Session *SystemSession `bson:"session"`
	}

	SystemSession struct {
		Issuer         string        `bson:"issuer"`
		ExpireAt       time.Duration `bson:"expireAt"`
		CheckIp        bool          `bson:"checkIp"`
		CheckUserAgent bool          `bson:"checkUserAgent"`
	}
)

func (x *DocSystem) GetID() primitive.ObjectID {
	return x.Id
}

func (x *DocSystem) GetColNm() string {
	return docSystem
}

func (x *DocSystem) GetMigrateList() mMigrate.FnMigrateList {
	return mgSystem
}

var mgSystem = mMigrate.FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "create init system"
		var model = mango.NewMDoc[System](docSystem)
		model.Data = &System{
			Session: &SystemSession{
				Issuer:         GetJwtIssuer(ctx),
				ExpireAt:       -1,
				CheckIp:        false,
				CheckUserAgent: true,
			},
		}
		err = model.Save(ctx)
		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */

type IReadSystem struct {
}

func (x *IReadSystem) Filter() (res bson.M, _ error) {
	res = bson.M{
		"_id": docSystemId,
	}
	return
}

func (x *IReadSystem) ColNm() string {
	return docSystem
}

func ReadSystem(ctx context.Context) (res *DocSystem, err error) {
	var system *mango.MDoc[System]
	if system, err = mango.ReadOneM[System](ctx, &IReadSystem{}); err != nil {
		return
	}
	res = fnReflect.ToPointer(DocSystem(*system))
	return
}

/* ------------------------------------------------------------------------------------------------------------ */

const ctxDocSystem = "CTX_DOC_SYSTEM"

var SetDocSystem = fnCtx.SetFn[*DocSystem](ctxDocSystem)
var GetDocSystem = fnCtx.GetFnP[*DocSystem](ctxDocSystem)

/* ------------------------------------------------------------------------------------------------------------ */

const ctxJwtSecret = "CTX_JWT_SECRET"

var SetJwtSecret = fnCtx.SetFn[string](ctxJwtSecret)
var GetJwtSecret = fnCtx.GetFnP[string](ctxJwtSecret)

/* ------------------------------------------------------------------------------------------------------------ */

const ctxJwtIssuer = "CTX_JWT_ISSUER"

var SetJwtIssuer = fnCtx.SetFn[string](ctxJwtIssuer)
var GetJwtIssuer = fnCtx.GetFnP[string](ctxJwtIssuer)
