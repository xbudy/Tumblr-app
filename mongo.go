package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//InitMongoDb Return a mongodb client
func InitMongoDb() (*mongo.Client, context.Context) {
	dbSrv := "mongodb://127.0.0.1:27017/myFirstDatabase?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(dbSrv))
	check_err(err)
	ctx := context.Background()
	err = client.Connect(ctx)
	check_err(err)
	return client, ctx
}

func (Mdb MongoDb) GetData(blog string) []post {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	cursor, err := col.Find(Mdb.Ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	var data []post
	cursor.All(Mdb.Ctx, &data)
	return data
}
func check_err(err error) {
	if err != nil {
		log.Error(err)
	}
}

type post struct {
	Images    []string
	Timestamp int
	Id        string
}

func addError(Post post, blog string) {
	log.Error("err added")
	cl, ctx := InitMongoDb()
	db := cl.Database("tumblr")
	col := db.Collection(blog + "" + "Errors")
	col.InsertOne(ctx, Post)
	defer cl.Disconnect(ctx)

}
func GetErrors(blog string) []post {
	cl, ctx := InitMongoDb()
	db := cl.Database("tumblr")
	col := db.Collection(blog + "" + "Errors")
	cursor, _ := col.Find(ctx, bson.D{})
	var data []post
	cursor.All(ctx, &data)
	return data
}

func RemoveError(blog string, Post post) {
	cl, ctx := InitMongoDb()
	db := cl.Database("tumblr")
	col := db.Collection(blog + "" + "Errors")
	col.DeleteOne(ctx, Post)
}
func (Mdb MongoDb) getPost(blog, Id string) post {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	var Post post
	r := col.FindOne(Mdb.Ctx, bson.D{{Key: "Id", Value: Id}})
	r.Decode(&Post)
	return Post
}
func (Mdb MongoDb) GetPostsCount(blog string) (int, error) {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	count, err := col.CountDocuments(Mdb.Ctx, bson.D{})
	return int(count), err
}

type MongoDb struct {
	Client *mongo.Client
	Ctx    context.Context
}
