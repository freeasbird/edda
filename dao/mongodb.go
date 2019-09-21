package dao

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"io"
	"time"
)

type CursorFunc func(cursor *mongo.Cursor) (err error)

type MongoCli struct {
	options     *Options
	host        string //  47.94.99.171:27017
	client      *mongo.Client
	database    *mongo.Database
	Collections map[string]*mongo.Collection
	timeout     time.Duration
}

type CursorF func(*mongo.Cursor) error

func (m *MongoCli) Init(ctx context.Context, opts ...Option) (err error) {
	m.options = new(Options)
	for _, opt := range opts {
		opt(m.options)
	}
	m.host = m.options.Host + ":" + m.options.Port
	m.timeout = m.options.Timeout
	//m.client, err =mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+username+":"+password+"@"+m.host))
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+m.host), options.Client().SetAuth(options.Credential{Username: m.options.Username, Password: m.options.Password, AuthSource: m.options.Database}))
	if err != nil {
		log.Fatal(err)
		return
	}
	m.database = m.client.Database(m.options.Database)
	// 2，选择数据库 my_db 表
	m.Collections = make(map[string]*mongo.Collection)
	for coll, index := range m.options.CollIndex {
		m.Collections[coll] = m.database.Collection(coll)
		// 唯一索引
		if index != "" {
			indexModel := mongo.IndexModel{
				Keys:    bsonx.Doc{{index, bsonx.Int32(1)},},
				Options: options.Index().SetUnique(true),
			}
			//indexModel = mongo.IndexModel{
			//	Keys: bsonx.Doc{{"expire_date", bsonx.Int32(1)}}, // 设置TTL索引列"expire_date"
			//	Options:options.Index().SetExpireAfterSeconds(1*24*3600), // 设置过期时间1天，即，条目过期一天过自动删除
			//}
			_, err = m.Collections[coll].Indexes().CreateOne(
				ctx,
				indexModel,
				options.CreateIndexes(),
			)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	// 3,测试连接
	err = m.client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// 参数 数据集名称  数据
func (m *MongoCli) Insert(coll string, data interface{}) (id string, err error) {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	result, err := Coll.InsertOne(ctx, data)
	if err != nil {
		log.Errorf("insertOne error.err:%s", err.Error())
		return
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", err
	}
	return oid.Hex(), err
}

func (m *MongoCli) Delete(coll string, filter interface{}) (err error) {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	_, err = Coll.DeleteOne(ctx, filter)
	return
}

func (m *MongoCli) FindOne(coll string, filter interface{}, result interface{}) (err error) {
	var (
		single *mongo.SingleResult
	)
	Coll, ok := m.Collections[coll]
	if !ok {
		return errors.New("page not found.")
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	if single = Coll.FindOne(ctx, filter); err != nil {
		log.Error("find error.")
		return
	}

	if err = single.Decode(result); err != nil {
		return
	}
	return
}

func (m *MongoCli) Exist(coll string, filter interface{}) bool {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return false
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	if count, err := Coll.CountDocuments(ctx, filter, options.Count()); err != nil {
		return false
	} else {
		return count > 0
	}
}

func (m *MongoCli) Find(coll string, filter interface{}, cursorF interface{}, skip, limit int64, sort int) (err error) {
	var (
		cursor *mongo.Cursor
	)
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	opt := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"_id": sort})
	if cursor, err = Coll.Find(ctx, filter, opt); err != nil {
		log.Debug("find error.")
		return
	}

	// 释放游标
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Error("cursor.Close error.", err)
		}
	}()

	f, ok := cursorF.(func(*mongo.Cursor) error)
	if ok {
		return f(cursor)
	}
	return errors.New("func type error.")
}

func (m *MongoCli) Count(coll string, filter interface{}) (num int64, err error) {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	return Coll.CountDocuments(ctx, filter)
}

func (m *MongoCli) FindOneUpdate(coll string, filter, update interface{}) (err error) {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	result := Coll.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetSort(bson.M{"_id": -1}))
	return result.Err()
}

func (m *MongoCli) Last(coll string) (last interface{}, err error) {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	err = Coll.FindOne(ctx, bson.D{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(last)
	return last, err
}

// 聚合查询
func (m *MongoCli) Aggregation(coll string, pipe interface{}, cursorF interface{}) (err error) {
	var (
		cursor *mongo.Cursor
	)
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	opt := options.Aggregate()
	if cursor, err = Coll.Aggregate(ctx, pipe, opt); err != nil {
		log.Debug("find error.")
		return
	}

	// 释放游标
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Error("cursor.Close error.", err)
		}
	}()

	f, ok := cursorF.(func(*mongo.Cursor) error)
	if ok {
		return f(cursor)
	}
	return errors.New("func type error.")

}

func (m *MongoCli) LastUpdate(coll string, update interface{}) (err error) {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	//err=m.Collection.FindOne(ctx,bson.D{},options.FindOne().SetSort(bson.M{"_id": -1})).Decode(data)
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	result := Coll.FindOneAndUpdate(ctx, bson.D{}, update, options.FindOneAndUpdate().SetSort(bson.M{"_id": -1}))
	return result.Err()
}

func (m *MongoCli) Update(coll string, filter interface{}, update interface{}) (err error) {
	Coll, ok := m.Collections[coll]
	if !ok {
		log.Error("Collection is not exist.")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	//Coll.FindOne(ctx,bson.D{},options.FindOne().SetSort(bson.M{"_id": -1})).Decode(data)
	_, err = Coll.UpdateOne(ctx, filter, update)
	return
}

func (m *MongoCli) Upload(name string, source io.Reader) (id string, err error) {
	bucket, err := gridfs.NewBucket(m.database)
	if err != nil {
		return
	}
	objId, err := bucket.UploadFromStream(name, source, )
	if err != nil {
		return
	}
	return objId.String(), err

}

func (m *MongoCli) Download(id interface{}, stream io.Writer) (size int64, err error) {
	bucket, err := gridfs.NewBucket(m.database)
	if err != nil {
		return
	}
	size, err = bucket.DownloadToStream(id, stream)
	return
}

func (m *MongoCli) FindFile(filter interface{}, cursorF interface{}, skip, limit int32) (err error) {
	bucket, err := gridfs.NewBucket(m.database)
	if err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)
	opt := options.GridFSFind().SetSkip(skip).SetLimit(limit)
	cursor, err := bucket.Find(filter, opt)
	if err != nil {
		return
	}
	// 释放游标
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Error("cursor.Close error.", err)
		}
	}()

	f, ok := cursorF.(func(*mongo.Cursor) error)
	if ok {
		return f(cursor)
	}
	return errors.New("func type error.")
}

func (m *MongoCli) DeleteFile(id interface{}) (err error) {
	bucket, err := gridfs.NewBucket(m.database)
	if err != nil {
		return
	}
	return bucket.Delete(id)
}
