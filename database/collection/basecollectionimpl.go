package collection

import (
	"context"
	"fmt"
	"go-gin-restful-service/constants"
	"go-gin-restful-service/log"
	"time"

	"github.com/bwmarrin/snowflake"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseCollectionImpl
// @Description: collection 的实现
type BaseCollectionImpl struct {
	DbName     string
	ColName    string
	DataBase   *mongo.Database
	Collection *mongo.Collection
}

var loc, _ = time.LoadLocation("Asia/Shanghai")

func AppendAutoCondition(filter bson.D, tenantId *string, delFlag int) bson.D {
	filter = append(filter, bson.E{
		Key:   "delFlag",
		Value: delFlag,
	})

	return filter
}

func GenerateSnowID() string {
	time.Sleep(1 * time.Nanosecond)
	snowflake.StepBits = 16
	node, err := snowflake.NewNode(999)
	if err != nil {
		log.Logger.Panic(err)
		return "0"
	}
	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	log.Logger.Infof("String ID: %s\n", id)
	return id.String()
}

func autoFillCreateFields(data interface{}) map[string]interface{} {
	myMap, ok := data.(map[string]interface{})
	if !ok {
		myMap = data.(primitive.M)
	}
	// 默认添加_id
	myMap["_id"] = GenerateSnowID()
	myMap["id"] = myMap["_id"]
	myMap["createTime"] = time.Now().In(loc).Format(constants.DateTimeFormatString)
	myMap["updateTime"] = time.Now().In(loc).Format(constants.DateTimeFormatString)
	myMap["version"] = 1
	myMap["delFlag"] = 0
	return myMap
}

func autoFillModifiedFields(data interface{}) map[string]interface{} {
	myMap, ok := data.(map[string]interface{})
	if !ok {
		myMap = data.(primitive.M)
	}
	myMap["updateTime"] = time.Now().In(loc).Format(constants.DateTimeFormatString)
	if val, ok := myMap["version"]; ok {
		myMap["version"] = val.(int64) + 1
	}
	return myMap
}

func (b *BaseCollectionImpl) SelectPage(ctx context.Context, filter interface{}, order *string, sort *int, skip, limit int64, excludeFileds *map[string]int) (int64, []map[string]interface{}, error) {
	var err error

	resultCount, err := b.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}
	if order == nil || len(*order) == 0 {
		u := string("updateTime")
		order = &u
	}
	if sort == nil {
		s := int(-1)
		sort = &s
	}
	defaultExcludes := map[string]int{
		"delFlag": 0,
		"version": 0,
		"_id":     0,
	}
	if excludeFileds != nil {
		for k, v := range *excludeFileds {
			_, isMapContainsKey := defaultExcludes[k]
			if v > 0 && isMapContainsKey { // Let the caller has the ability to remove defaultExclude
				delete(defaultExcludes, k)
			} else {
				defaultExcludes[k] = v
			}
		}
	}
	opts := options.Find().SetSort(bson.D{{Key: *order, Value: *sort}}).SetSkip(skip).SetLimit(limit).SetProjection(defaultExcludes)
	finder, err := b.Collection.Find(ctx, filter, opts)
	if err != nil {
		return resultCount, nil, err
	}

	result := make([]map[string]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return resultCount, nil, err
	}
	return resultCount, result, nil
}

func (b *BaseCollectionImpl) SelectList(ctx context.Context, filter interface{}, order *string, sort *int, excludeFileds *map[string]int) ([]map[string]interface{}, error) {
	var err error
	if order == nil {
		u := string("updateTime")
		order = &u
	}
	if sort == nil {
		s := int(-1)
		sort = &s
	}
	defaultExcludes := map[string]int{
		"delFlag": 0,
		"version": 0,
		"_id":     0,
	}
	if excludeFileds != nil {
		for k, v := range *excludeFileds {
			_, isMapContainsKey := defaultExcludes[k]
			if v > 0 && isMapContainsKey { // Let the caller has the ability to remove defaultExclude
				delete(defaultExcludes, k)
			} else {
				defaultExcludes[k] = v
			}
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: *order, Value: *sort}}).SetProjection(defaultExcludes)
	finder, err := b.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (b *BaseCollectionImpl) SelectOne(ctx context.Context, filter interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	err := b.Collection.FindOne(ctx, filter, options.FindOne()).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *BaseCollectionImpl) SelectCount(ctx context.Context, filter interface{}) (int64, error) {
	return b.Collection.CountDocuments(ctx, filter)
}

func (b *BaseCollectionImpl) UpdateOne(ctx context.Context, filter, update interface{}) (int64, error) {
	update = autoFillModifiedFields(update)
	update = bson.M{"$set": update}
	result, err := b.Collection.UpdateOne(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

func (b *BaseCollectionImpl) UpdateMany(ctx context.Context, filter, update interface{}) (int64, error) {
	update = autoFillModifiedFields(update)
	update = bson.M{"$set": update}
	result, err := b.Collection.UpdateMany(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}
func (b *BaseCollectionImpl) LogicDelete(ctx context.Context, filter interface{}) (int64, error) {
	update := bson.M{"delFlag": constants.DELFLAG_DELETED}
	update = autoFillModifiedFields(update)
	update = bson.M{"$set": update}
	result, err := b.Collection.UpdateMany(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Logic Delete result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

func (b *BaseCollectionImpl) RollbackDeleted(ctx context.Context, filter interface{}) (int64, error) {
	update := bson.M{"delFlag": constants.DELFLAG_NORMAL}
	update = autoFillModifiedFields(update)
	update = bson.M{"$set": update}
	result, err := b.Collection.UpdateMany(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Rollback Deleted result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

func (b *BaseCollectionImpl) Delete(ctx context.Context, filter interface{}) (int64, error) {
	result, err := b.Collection.DeleteMany(ctx, filter, options.Delete())
	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("DeleteOne result: %s ", "document not found")
	}
	return result.DeletedCount, nil
}

func (b *BaseCollectionImpl) InsertOne(ctx context.Context, model interface{}) (interface{}, error) {
	//批量插入数据，自动填充创建人和更新人字段
	model = autoFillCreateFields(model)
	result, err := b.Collection.InsertOne(ctx, model, options.InsertOne())
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

func (b *BaseCollectionImpl) InsertMany(ctx context.Context, models []interface{}) ([]interface{}, error) {
	//批量插入数据，自动填充创建人和更新人字段
	tobeSaved := make([]interface{}, 0)
	for _, item := range models {
		tobeSaved = append(tobeSaved, autoFillCreateFields(item))
	}

	result, err := b.Collection.InsertMany(ctx, tobeSaved, options.InsertMany())
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, err
}

func (b *BaseCollectionImpl) Aggregate(ctx context.Context, pipeline interface{}, result interface{}) error {
	finder, err := b.Collection.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return err
	}
	if err := finder.All(ctx, &result); err != nil {
		return err
	}
	return nil
}

func (b *BaseCollectionImpl) CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error {
	_, err := b.Collection.Indexes().CreateMany(ctx, indexes, options.CreateIndexes())
	return err
}

func (b *BaseCollectionImpl) GetCollection() *mongo.Collection {
	return b.Collection
}
