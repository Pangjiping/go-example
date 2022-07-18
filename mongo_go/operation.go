package mongo_go

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 因为mongo中数据是使用bson格式存储的，所以在定义插入数据结构体时，需要使用bson标签进行反射

// 任务执行的时间点
type TimePoint struct {
	StartTime int64 `bson:"start_time"`
	EndTime   int64 `bson:"end_time"`
}

// 日志结构
type LogRecord struct {
	JobName   string    `bson:"job_name"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"time_point"`
}

// InsertOne2Mongo 插入数据
func InsertOne2Mongo() error {
	record := &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	// 插入数据
	collection, err := getCollection()
	if err != nil {
		return err
	}

	res, err := collection.InsertOne(context.Background(), record)
	if err != nil {
		return err
	}

	// 默认生成一个全局唯一的id，12字节的二进制
	docID := res.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID: ", docID.Hex())

	return nil
}

func InsertMany2Mongo() error {
	record := &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	// 声明一个插入数据的切片
	logArr := []interface{}{record, record, record}

	collection, err := getCollection()
	if err != nil {
		return err
	}

	result, err := collection.InsertMany(context.Background(), logArr)
	if err != nil {
		return err
	}

	for _, id := range result.InsertedIDs {
		docID := id.(primitive.ObjectID)
		fmt.Println("自增ID: ", docID.Hex())
	}

	return nil
}

// jobname过滤条件
type FindByJobName struct {
	JobName string `bson:"job_name"`
}

func FindFromMongo() error {

	collection, err := getCollection()
	if err != nil {
		return err
	}

	// 按照jobname过滤，找出jobname=job10
	cond := &FindByJobName{
		JobName: "job10",
	}
	skip := int64(0)
	limit := int64(2)

	// 查询
	//cursor, err := collection.Find(context.Background(), bson.D{{"jon_name", "job10"}}) // 这是官方案例中的写法，显得不像在写客户端

	cursor, err := collection.Find(context.Background(), cond, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	// 遍历结果集
	for cursor.Next(context.Background()) {
		record := &LogRecord{}

		// 反序列化bson到结构体对象
		err := cursor.Decode(record)
		if err != nil {
			return err
		}

		fmt.Println(*record)
	}
	return nil
}

// 根据删除规则我们定义一个过滤器，我们要删除的是创建时间小于当前时间的，那么在mongo ctl中就应该写成json表达式
//
// delete({"timePoint.startPoint":{"$lt":当前时间}})
//
// 我们可以利用go的bson反射来做到这个表达式的定义，记住反射是怎么序列化的就行，后面的bson标签是key

// startTime小于某时间
// {"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// 定义删除条件
// {"timePoint.startPoint":{"$lt":timestamp}}
type DelCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func DeleteFromMongo() error {
	collection, err := getCollection()
	if err != nil {
		return err
	}

	// 删除条件
	delCond := &DelCond{
		beforeCond: TimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}

	result, err := collection.DeleteMany(context.Background(), delCond)
	if err != nil {
		return err
	}

	fmt.Println("删除了多少行: ", result.DeletedCount)
	return nil

}
