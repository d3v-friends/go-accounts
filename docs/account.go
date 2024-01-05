package docs

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mMigrate"
	"github.com/d3v-friends/mango/mType"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	docAccount             = "accounts"
	fAccountDataIdentifier = "data.identifier"
	fAccountDataProperties = "data.properties"
)

type (
	DocAccount  mango.MDoc[Account]
	IdKey       string
	PropKey     string
	VerifierKey string

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

// CreateAccount
// 유저 생성
func CreateAccount(ctx context.Context, i *ICreateAccount) (res *DocAccount, err error) {
	var account = mango.NewDoc[Account](docAccount, i.Data)

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
	filter = make(bson.M)

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
// init indexing

type IReindexingAccount struct {
	IdentifierKeys []string
	PropertyKeys   []string
}

func (x IReindexingAccount) ForEachIdentifier(fn func(fieldNm string, idxNm string) error) (err error) {
	for _, v := range x.IdentifierKeys {
		if err = fn(
			fmt.Sprintf("data.identifer.%s", v),
			fmt.Sprintf("data.identifier.%s_1", v),
		); err != nil {
			return
		}
	}
	return
}

func (x IReindexingAccount) ForEachProperty(fn func(fieldNm string, idxNm string) error) (err error) {
	for _, v := range x.PropertyKeys {
		if err = fn(
			fmt.Sprintf("data.property.%s", v),
			fmt.Sprintf("data.property.%s_1", v)); err != nil {
			return
		}
	}
	return
}

func (x IReindexingAccount) Has(key string) (has bool) {
	for _, v := range x.IdentifierKeys {
		if fmt.Sprintf("data.identifier.%s_1", v) == key {
			return true
		}
	}

	for _, v := range x.PropertyKeys {
		if fmt.Sprintf("data.property.%s_1", v) == key {
			return true
		}
	}

	return false
}

func ReindexingAccount(
	ctx context.Context,
	i *IReindexingAccount,
) (err error) {
	var col = mango.GetColP(ctx, docAccount)

	var cur *mongo.Cursor
	if cur, err = col.Indexes().List(ctx); err != nil {
		return
	}

	// 추가된 키 찾아서 인덱싱키 만들기
	var loadIdxes = make(mType.IndexModels, 0)
	if err = cur.All(ctx, &loadIdxes); err != nil {
		return
	}

	if err = i.ForEachIdentifier(func(fieldNm, idxNm string) (fnErr error) {
		if loadIdxes.Has(idxNm) {
			return
		}

		_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   fieldNm,
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Unique: fnReflect.ToPointer(true),
			},
		})

		return
	}); err != nil {
		return
	}

	if err = i.ForEachProperty(func(fieldNm, idxNm string) (fnErr error) {
		if loadIdxes.Has(idxNm) {
			return
		}

		_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   fieldNm,
					Value: 1,
				},
			},
		})

		return
	}); err != nil {
		return
	}

	// 없는것 삭제하기
	for _, v := range loadIdxes {
		if i.Has(v.Name) {
			continue
		}

		if _, err = col.Indexes().DropOne(ctx, v.Name); err != nil {
			return
		}
	}

	return
}
