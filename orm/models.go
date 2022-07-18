package orm

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id      int
	Name    string
	Profile *Profile `orm:"rel(one)"`      // OneToOne relation
	Post    []*Post  `orm:"reverse(many)"` // 设置一对多的反向关系
}

type Profile struct {
	Id   int
	Age  int16
	User *User `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}

type Post struct {
	Id    int
	Title string
	User  *User  `orm:"rel(fk)"` //设置一对多关系
	Tags  []*Tag `orm:"rel(m2m)"`
}

type Tag struct {
	Id    int
	Name  string
	Posts []*Post `orm:"reverse(many)"` //设置多对多反向关系
}

var o orm.Ormer

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// 注册数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@(123.57.33.149:3306)/orm_test?charset=utf8")

	// 注册模型，对于使用orm.QuerySeter进行高级查询是必须的
	orm.RegisterModel(new(User), new(Profile), new(Post), new(Tag))

	// 使用表名前缀和后缀
	//orm.RegisterModelWithPrefix("prefix_", new(User))
	//orm.RegisterModelWithSuffix("_suffix", new(User))

	// 自动建表
	orm.RunSyncdb("default", true, true)

	// 最大连接数
	orm.SetMaxOpenConns("default", 30)

	// 最大空闲连接
	orm.SetMaxIdleConns("default", 30)

	// 设置时区
	orm.DefaultTimeLoc = time.UTC

	// 初始化ormer
	o = orm.NewOrm()

	// 初始化数据
	// datainit()
}
