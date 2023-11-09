package docs

import (
	"context"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango/mMigrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const docAccount = "accounts"

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
						Key:   "data.identifiers.key",
						Value: 1,
					},
					{
						Key:   "data.identifier.value",
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
						Key:   "data.properties.key",
						Value: 1,
					},
					{
						Key:   "data.properties.value",
						Value: 1,
					},
				},
			},
		})
		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */

type ICreateAccount struct {
	Identifiers []*KV
	Verifiers   []*Verifier
	Properties  []*KV
}

func CreateAccount(ctx context.Context, i *ICreateAccount) (res *Account, err error) {
	panic("not impl")
}

/* ------------------------------------------------------------------------------------------------------------ */

type IReadAccount struct {
	Ids []*primitive.ObjectID
}
