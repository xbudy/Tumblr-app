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
func (Mdb MongoDb) GetPosts(blog string) []PostMeta {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	cursor, err := col.Find(Mdb.Ctx, bson.D{{Key: "type", Value: "post"}})
	if err != nil {
		panic(err)
	}
	var data []PostMeta
	cursor.All(Mdb.Ctx, &data)
	return data
}
func (Mdb MongoDb) GetData(blog string) []post {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	cursor, err := col.Find(Mdb.Ctx, bson.D{{Key: "type", Value: "post"}})
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
func (Mdb MongoDb) getPost(blog, Id string) PostMeta {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	var Post PostMeta
	r := col.FindOne(Mdb.Ctx, bson.D{{Key: "id", Value: Id}})
	r.Decode(&Post)
	if r.Err() != nil {
		log.Println(r.Err(), Id)
	}
	return Post
}
func (Mdb MongoDb) GetPostsCount(blog string) (int, error) {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	count, err := col.CountDocuments(Mdb.Ctx, bson.D{{Key: "type", Value: "post"}})
	return int(count), err
}
func (Mdb MongoDb) AddPost(blog string, post PostMeta) {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	//_, err := col.InsertOne(Mdb.Ctx, post)
	opts := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(Mdb.Ctx, post, bson.D{{Key: "$set", Value: post}}, opts)
	if err != nil {
		log.Println(err)
	}
}
func (Mdb MongoDb) setLastPage(blog string, lastPage string) {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	opts := options.Update().SetUpsert(true)
	_, e := col.UpdateOne(Mdb.Ctx, bson.D{{Key: "type", Value: "lastPage"}}, bson.D{{Key: "$set", Value: bson.D{{Key: "value", Value: lastPage}}}}, opts)
	if e != nil {
		log.Error(e)
	}
}
func (Mdb MongoDb) GetLastPage(blog string) string {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	res := col.FindOne(Mdb.Ctx, bson.D{{Key: "type", Value: "lastPage"}})
	var data lastPage
	res.Decode(&data)
	return data.Value
}
func (Mdb MongoDb) SetBlogInfo(blog string, binfo bloginfo) {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	opts := options.Update().SetUpsert(true)
	col.UpdateOne(Mdb.Ctx, bson.D{{Key: "type", Value: "bloginfo"}}, bson.D{{Key: "$set", Value: bson.D{{Key: "total_posts", Value: binfo.TotalPosts}}}}, opts)
}
func (Mdb MongoDb) GetBlogInfo(blog string) bloginfo {
	db := Mdb.Client.Database("tumblr")
	col := db.Collection(blog)
	res := col.FindOne(Mdb.Ctx, bson.D{{Key: "type", Value: "bloginfo"}})

	var data bloginfo
	res.Decode(&data)
	return data
}

type bloginfo struct {
	TotalPosts int `bson:"total_posts"`
}
type lastPage struct {
	Value string
}
type MongoDb struct {
	Client *mongo.Client
	Ctx    context.Context
}
