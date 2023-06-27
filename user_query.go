package go_mongo

import (
	"cc/go-mongo/user"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserQuery struct {
	config
	Predicates []user.UserPredicate
	limit      *int64
	offset     *int64
	dbName     string
	options    bson.D
}

func (uq *UserQuery) Limit(limit int64) *UserQuery {
	uq.limit = &limit
	return uq
}

func (uq *UserQuery) Offset(offset int64) *UserQuery {
	uq.offset = &offset
	return uq
}

func (uq *UserQuery) Order(o ...OrderFunc) *UserQuery {
	for _, fn := range o {
		fn(&uq.options)
	}
	return uq
}

func (uq *UserQuery) Where(ps ...user.UserPredicate) *UserQuery {
	for _, p := range ps {
		uq.Predicates = append(uq.Predicates, p)
	}
	return uq
}

func (uq *UserQuery) All(ctx context.Context) ([]*User, error) {
	filter := bson.D{}
	for _, p := range uq.Predicates {
		p(&filter)
	}

	o := options.Find()
	if uq.limit != nil && *uq.limit != 0 {
		o = o.SetLimit(*uq.limit)
	}
	if uq.offset != nil && *uq.offset != 0 {
		o = o.SetSkip(*uq.offset)
	}
	o.SetSort(uq.options)
	cur, err := uq.Database(uq.dbName).Collection(user.UserMongo).Find(ctx, filter, o)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	temp := make([]*User, 0, 0)
	for cur.Next(ctx) {
		var u User
		err = cur.Decode(&u)
		if err != nil {
			return nil, err
		}
		temp = append(temp, &u)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}
	return temp, nil
}
