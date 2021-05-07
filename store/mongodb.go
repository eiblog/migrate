// Package store provides ...
package store

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/eiblog/migrate/v1"
	v2 "github.com/eiblog/migrate/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongodb struct{}

// LoadEiBlog 读取数据
func (db mongodb) LoadEiBlog(from Store) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	opts := options.Client().ApplyURI(from.Source)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	switch from.Version {
	case "v1":
		// account & blogger & serie
		collection := client.Database(v1.DB).Collection(v1.COLLECTION_ACCOUNT)
		filter := bson.M{}
		acct := v1.Account{}
		err = collection.FindOne(ctx, filter).Decode(&acct)
		if err != nil {
			return nil, err
		}
		// articles
		collection = client.Database(v1.DB).Collection(v1.COLLECTION_ARTICLE)
		filter = bson.M{}
		var articles []v1.Article
		cur, err := collection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
		for cur.Next(ctx) {
			obj := v1.Article{}
			err = cur.Decode(&obj)
			if err != nil {
				return nil, err
			}
			articles = append(articles, obj)
		}
		cur.Close(ctx)

		blog := &v1.EiBlog{
			Account:  acct,
			Articles: articles,
		}
		return blog, nil
	case "v2":
		// blogger
		collection := client.Database(v2.MongoDBName).Collection(v2.CollectionBlogger)
		filter := bson.M{}
		blogger := v2.Blogger{}
		err = collection.FindOne(ctx, filter).Decode(&blogger)
		if err != nil {
			return nil, err
		}
		// account
		collection = client.Database(v2.MongoDBName).Collection(v2.CollectionAccount)
		filter = bson.M{}
		acct := v2.Account{}
		err = collection.FindOne(ctx, filter).Decode(&acct)
		if err != nil {
			return nil, err
		}
		// articles
		collection = client.Database(v2.MongoDBName).Collection(v2.CollectionArticle)
		filter = bson.M{}
		var articles []v2.Article
		cur, err := collection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
		for cur.Next(ctx) {
			obj := v2.Article{}
			err = cur.Decode(&obj)
			if err != nil {
				return nil, err
			}
			articles = append(articles, obj)
		}
		cur.Close(ctx)
		// serie
		collection = client.Database(v2.MongoDBName).Collection(v2.CollectionSerie)
		filter = bson.M{}
		var series []v2.Serie
		cur, err = collection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
		for cur.Next(ctx) {
			obj := v2.Serie{}
			err = cur.Decode(&obj)
			if err != nil {
				return nil, err
			}
			series = append(series, obj)
		}
		cur.Close(ctx)
		blog := &v2.EiBlog{
			Blogger:  blogger,
			Account:  acct,
			Articles: articles,
			Series:   series,
		}
		return blog, nil
	}
	return nil, fmt.Errorf("unsupported version: %s", from.Version)
}

// StoreEiBlog 保存数据
func (db mongodb) StoreEiBlog(to Store, blog interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(to.Source)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	switch data := blog.(type) {
	case *v1.EiBlog:
		// account
		collection := client.Database(v1.DB).Collection(v1.COLLECTION_ACCOUNT)
		_, err = collection.InsertOne(ctx, data.Account)
		if err != nil {
			return err
		}
		// articles
		collection = client.Database(v1.DB).Collection(v1.COLLECTION_ARTICLE)
		var maxArticleID int32 = 0
		for _, v := range data.Articles {
			_, err = collection.InsertOne(ctx, v)
			if err != nil {
				return err
			}
			if v.ID > maxArticleID {
				maxArticleID = v.ID
			}
		}
		// serie id
		var maxSerieID int32 = 0
		for _, v := range data.Account.Series {
			if v.ID > maxSerieID {
				maxSerieID = v.ID
			}
		}
		// counter
		collection = client.Database(v1.DB).Collection("COUNTERS")
		counter := bson.M{"name": v1.COUNTER_ARTICLE, "nextval": maxArticleID}
		_, err = collection.InsertOne(ctx, counter)
		if err != nil {
			return err
		}
		counter = bson.M{"name": v1.COUNTER_SERIE, "nextval": maxSerieID}
		_, err = collection.InsertOne(ctx, counter)
		if err != nil {
			return err
		}
	case *v2.EiBlog:
		// blogger
		collection := client.Database(v2.MongoDBName).Collection(v2.CollectionBlogger)
		_, err = collection.InsertOne(ctx, data.Blogger)
		if err != nil {
			return err
		}
		// account
		collection = client.Database(v2.MongoDBName).Collection(v2.CollectionAccount)
		_, err = collection.InsertOne(ctx, data.Account)
		if err != nil {
			return err
		}
		// articles
		collection = client.Database(v2.MongoDBName).Collection(v2.CollectionArticle)
		maxArticleID := 0
		for _, v := range data.Articles {
			_, err = collection.InsertOne(ctx, v)
			if err != nil {
				return err
			}
			if v.ID > maxArticleID {
				maxArticleID = v.ID
			}
		}
		// series
		collection = client.Database(v2.MongoDBName).Collection(v2.CollectionSerie)
		maxSerieID := 0
		for _, v := range data.Series {
			_, err = collection.InsertOne(ctx, v)
			if err != nil {
				return err
			}
			if v.ID > maxSerieID {
				maxSerieID = v.ID
			}
		}
		// counter
		collection = client.Database(v2.MongoDBName).Collection(v2.CollectionCounter)
		counter := bson.M{"name": v2.CounterNameArticle, "nextval": maxArticleID}
		_, err = collection.InsertOne(ctx, counter)
		if err != nil {
			return err
		}
		counter = bson.M{"name": v2.CounterNameSerie, "nextval": maxSerieID}
		_, err = collection.InsertOne(ctx, counter)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported data version: %T", blog)
	}
	return nil
}
