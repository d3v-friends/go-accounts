package docs

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-accounts/fn/fnJWT"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mMigrate"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	docSession = "sessions"
)

type (
	DocSession struct {
		Id         primitive.ObjectID `bson:"_id"`
		AccountId  primitive.ObjectID `bson:"accountId"`
		IsActivate bool               `bson:"isActivate"`
		IP         string             `bson:"ip"`
		UserAgent  string             `bson:"userAgent"`
		Nm         string             `bson:"nm"`
		SignInAt   time.Time          `bson:"signInAt"`
		LastConnAt time.Time          `bson:"lastConnAt"`
	}

	DocSessionAccount struct {
		*DocSession `bson:"inline"`
		Account     *DocAccount `bson:"account"`
	}
)

func (x DocSession) GetID() primitive.ObjectID {
	return x.Id
}

func (x DocSession) GetColNm() string {
	return docSession
}

func (x DocSession) GetMigrateList() mMigrate.FnMigrateList {
	return mgSession
}

func (x DocSession) GetSessionID() string {
	return x.Id.Hex()
}

func (x DocSession) GetAudience() []string {
	return []string{x.AccountId.Hex()}
}

func (x DocSession) GetSubject() string {
	return "sign_in"
}

/* ------------------------------------------------------------------------------------------------------------ */

var mgSession = mMigrate.FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{
						Key:   "_id",
						Value: -1,
					},
					{
						Key:   "isActivate",
						Value: 1,
					},
				},
				Options: &options.IndexOptions{
					Unique: fnReflect.ToPointer(true),
				},
			},
			{
				Keys: bson.D{
					{
						Key:   "accountId",
						Value: 1,
					},
				},
			},
			{
				Keys: bson.D{
					{
						Key:   "ip",
						Value: 1,
					},
				},
			},
			{
				Keys: bson.D{
					{
						Key:   "signInAt",
						Value: -1,
					},
				},
			},
			{
				Keys: bson.D{
					{
						Key:   "lastConnAt",
						Value: -1,
					},
				},
			},
		})

		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */

type ICreateSession struct {
	AccountId primitive.ObjectID
	Ip        string
	UserAgent string
}

func CreateSession(
	ctx context.Context,
	i *ICreateSession,
) (
	doc *DocSession,
	token string,
	err error,
) {
	var now = time.Now()
	doc = &DocSession{
		Id:         primitive.NewObjectID(),
		AccountId:  i.AccountId,
		IsActivate: true,
		IP:         i.Ip,
		UserAgent:  i.UserAgent,
		Nm:         i.UserAgent,
		SignInAt:   now,
		LastConnAt: now,
	}

	if _, err = mango.GetColP(ctx, docSession).InsertOne(ctx, doc); err != nil {
		return
	}

	var docSystem = GetDocSystem(ctx)
	var jwtSecret = GetJwtSecret(ctx)
	var jwt = fnJWT.NewV1[DocSession](
		jwtSecret,
		docSystem.Data.Session.Issuer,
		docSystem.Data.Session.ExpireAt,
	)

	if token, err = jwt.Encode(*doc); err != nil {
		return
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

type IVerifySession struct {
	Token     string
	Ip        string
	UserAgent string
}

func VerifySession(
	ctx context.Context,
	i *IVerifySession,
) (account *DocAccount, err error) {
	var docSystem = GetDocSystem(ctx)
	var jwtSecret = GetJwtSecret(ctx)
	var v1 = fnJWT.NewV1[DocSession](
		jwtSecret,
		docSystem.Data.Session.Issuer,
		docSystem.Data.Session.ExpireAt,
	)

	var claim *jwt.RegisteredClaims
	if claim, err = v1.Decode(i.Token); err != nil {
		return
	}

	var sessionId primitive.ObjectID
	if sessionId, err = primitive.ObjectIDFromHex(claim.ID); err != nil {
		return
	}

	var filter = &IReadSession{
		SessionId:  sessionId,
		IsActivate: true,
	}
	var session *DocSession
	if session, err = ReadOneSession(
		ctx,
		filter,
	); err != nil {
		return
	}

	if docSystem.Data.Session.CheckIp && session.IP != i.Ip {
		err = fmt.Errorf("invalid ip with session")
		return
	}

	if docSystem.Data.Session.CheckUserAgent && session.UserAgent != i.UserAgent {
		err = fmt.Errorf("invalid user_agent with session")
		return
	}

	if session, err = UpdateOneSession(ctx, filter, &IUpdateSession{
		Ip:        &i.Ip,
		UserAgent: &i.UserAgent,
	}); err != nil {
		return
	}

	return ReadOneAccount(ctx, &IReadAccount{
		Ids:        []primitive.ObjectID{session.AccountId},
		IsActivate: fnReflect.ToPointer(true),
	})
}

/* ------------------------------------------------------------------------------------------------------------ */

type IReadSession struct {
	SessionId  primitive.ObjectID
	IsActivate bool
}

func (x IReadSession) Filter() (res bson.M, _ error) {
	res = bson.M{
		"_id":        x.SessionId,
		"isActivate": x.IsActivate,
	}
	return
}

func (x IReadSession) ColNm() string {
	return docSession
}

type IUpdateSession struct {
	IsActivate *bool
	Ip         *string
	UserAgent  *string
}

func (x *IUpdateSession) Update() (update bson.M, err error) {
	var now = time.Now()

	var set = bson.M{
		"lastConnAt": now,
	}

	if x.Ip != nil {
		set["ip"] = *x.Ip
	}

	if x.IsActivate != nil {
		set["isActivate"] = *x.IsActivate
	}

	if x.UserAgent != nil {
		set["userAgent"] = *x.UserAgent
	}

	update = bson.M{
		"$set": set,
	}

	return
}

func ReadOneSession(
	ctx context.Context,
	i mango.IfFilter,
) (doc *DocSession, err error) {
	return mango.ReadOne[DocSession](ctx, i)
}

func UpdateOneSession(
	ctx context.Context,
	i mango.IfFilter,
	u mango.IfUpdate,
) (doc *DocSession, err error) {
	return mango.UpdateOne[DocSession](ctx, i, u)
}
