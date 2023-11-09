package docs

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mMigrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	docAccount             = "accounts"
	fAccountDataIdentifier = "data.identifier"
	fAccountDataProperties = "data.properties"
)

type (
	Account struct {
		Id        primitive.ObjectID `bson:"_id"`
		Data      *AccountData       `bson:"data"`
		History   []*AccountData     `bson:"history"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}

	AccountData struct {
		Identifiers []*KV       `bson:"identifiers"`
		Verifiers   []*Verifier `bson:"verifiers"`
		Properties  []*KV       `bson:"properties"`
		CreatedAt   time.Time   `bson:"createdAt"`
	}

	Verifier struct {
		Kind  string `bson:"kind"`
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}
)

func (x Account) GetID() primitive.ObjectID {
	return x.Id
}

func (x Account) GetColNm() string {
	return docAccount
}

func (x Account) GetMigrateList() mMigrate.FnMigrateList {
	return mgAccount
}

var mgAccount = mMigrate.FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{
						Key:   "data.identifiers",
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
						Key:   "data.properties",
						Value: 1,
					},
				},
			},
		})
		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */
// Creator

type ICreateAccount struct {
	Identifiers []*KV
	Verifiers   []*Verifier
	Properties  []*KV
}

func CreateAccount(ctx context.Context, i *ICreateAccount) (res *Account, err error) {
	var now = time.Now()

	res = &Account{
		Id: primitive.NewObjectID(),
		Data: &AccountData{
			Identifiers: i.Identifiers,
			Verifiers:   i.Verifiers,
			Properties:  i.Properties,
			CreatedAt:   now,
		},
		History:   make([]*AccountData, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}

	var col = mango.GetColP(ctx, docAccount)

	if _, err = col.InsertOne(ctx, res); err != nil {
		return
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */
// Reader

func ReadOneAccount(
	ctx context.Context,
	i mango.IfFilter,
	opts ...*options.FindOneOptions,
) (*Account, error) {
	return mango.ReadOne[Account](ctx, i, opts...)
}

func ReadAllAccount(
	ctx context.Context,
	i mango.IfFilter,
	opts ...*options.FindOptions,
) ([]*Account, error) {
	return mango.ReadAll[Account](ctx, i, opts...)
}

func ReadListAccount(
	ctx context.Context,
	i mango.IfFilter,
	p mango.IfPager,
	opts ...*options.FindOptions,
) (ls []*Account, total int64, err error) {
	return mango.ReadList[Account](ctx, i, p, opts...)
}

/* ------------------------------------------------------------------------------------------------------------ */
// Filters

// AccountFilter
// _id 는 $in operator 로 검색한다
// properties 는 $and 로 검색한다
type AccountFilter struct {
	Ids        []primitive.ObjectID
	Identifier *KV
	Properties []KV
}

func (x AccountFilter) Filter() (filter bson.M, err error) {
	filter = make(bson.M)
	if len(x.Ids) != 0 {
		filter[fId] = bson.M{
			"$in": x.Ids,
		}
	}

	if x.Identifier != nil {
		filter[fAccountDataIdentifier] = bson.M{
			"$elemMatch": *x.Identifier,
		}
	}

	if len(x.Properties) != 0 {
		filter[fAccountDataProperties] = bson.M{
			"$all": x.Properties,
		}
	}

	if len(filter) == 0 {
		err = fmt.Errorf("empty account_filter")
		return
	}

	return
}

func (x AccountFilter) ColNm() string {
	return docAccount
}
