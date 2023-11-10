package docs

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mMigrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	docAccount             = "accounts"
	fAccountDataIdentifier = "data.identifier"
	fAccountDataProperties = "data.properties"
)

type (
	DocAccount mango.MDoc[Account]

	Account struct {
		IsActivate bool                        `bson:"isActivate"`
		Identifier map[string]string           `bson:"identifier"`
		Property   map[string]string           `bson:"property"`
		Verifier   map[string]*AccountVerifier `bson:"verifier"`
		Data       []byte                      `bson:"data"`
	}

	AccountVerifier struct {
		Key   string     `bson:"key"`
		Value string     `bson:"value"`
		Mode  VerifyMode `bson:"mode"`
	}
)

func (x *DocAccount) GetID() primitive.ObjectID {
	return x.Id
}

func (x *DocAccount) GetColNm() string {
	return docAccount
}

func (x *DocAccount) GetMigrateList() mMigrate.FnMigrateList {
	return make(mMigrate.FnMigrateList, 0)
}

/* ------------------------------------------------------------------------------------------------------------ */

type ICreateAccount struct {
	Data *Account
}

func CreateAccount(ctx context.Context, i *ICreateAccount) (res *DocAccount, err error) {
	var account = mango.NewMDoc[Account](docAccount)
	account.Data = i.Data

	if err = account.Save(ctx); err != nil {
		return
	}

	res = fnReflect.ToPointer(DocAccount(*account))
	return
}

/* ------------------------------------------------------------------------------------------------------------ */

// IReadAccount
// ids -> $in
// identifier -> exactly
// property -> like
type IReadAccount struct {
	Ids        []primitive.ObjectID
	IsActivate *bool
	Identifier map[string]string
	Property   map[string]string
}

func (x IReadAccount) Filter() (filter bson.M, err error) {
	if len(x.Ids) != 0 {
		filter["_id"] = bson.M{
			"$in": x.Ids,
		}
	}

	if len(x.Identifier) != 0 {
		for key, value := range x.Identifier {
			filter[fmt.Sprintf("data.identifiers.%s", key)] = value
		}
	}

	if len(x.Property) != 0 {
		for key, value := range x.Property {
			filter[fmt.Sprintf("data.properties.%s", key)] = bson.E{
				Key: "$regex",
				Value: primitive.Regex{
					Pattern: value,
				},
			}
		}
	}

	if x.IsActivate != nil {
		filter["data.isActivate"] = *x.IsActivate
	}

	return
}

func (x IReadAccount) ColNm() string {
	return docAccount
}

func ReadOneAccount(
	ctx context.Context,
	i mango.IfFilter,
	opts ...*options.FindOneOptions,
) (account *DocAccount, err error) {
	var res *mango.MDoc[Account]
	if res, err = mango.ReadOneM[Account](ctx, i, opts...); err != nil {
		return
	}
	account = fnReflect.ToPointer(DocAccount(*res))
	return
}

/* ------------------------------------------------------------------------------------------------------------ */
