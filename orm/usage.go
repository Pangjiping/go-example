package orm

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

//func init() {
//	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/db_test?charset=utf8&loc=Local")
//	orm.RegisterModel(new(User))
//	orm.RunSyncdb("default", false, true)
//}

// FastUsage orm的快速使用案例
func FastUsage() error {
	user := User{
		Name: "slene",
	}

	// insert
	id, err := o.Insert(&user)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %d, ERR: %v", id, err)

	// update
	user.Name = "astaxie"
	num, err := o.Update(&user)
	if err != nil {
		return err
	}
	fmt.Printf("NUM: %d, ERR: %v", num, err)

	// read
	u := User{Id: user.Id}
	err = o.Read(&u)
	if err != nil {
		return err
	}
	fmt.Printf("User: %v, ERR: %v", u, err)

	// delete
	num, err = o.Delete(&u)
	if err != nil {
		return err
	}
	fmt.Printf("NUM: %d, ERR: %v", num, err)

	return nil
}

// UnionSearch 关联查询
func UnionSearch() error {
	var posts []*Post

	qs := o.QueryTable("post")
	num, err := qs.Filter("User__Name", "slene").All(&posts)
	if err != nil {
		return err
	}
	fmt.Printf("NUM: %d, ERR: %v", num, err)

	return nil
}

// SQLSearch sql查询
// 当无法使用ORM来达到需求时，也可以直接使用SQL来完成查询、映射操作
func SQLSearch() error {

	var maps []orm.Params
	num, err := o.Raw("SELECT * FROM user").Values(&maps)
	if err != nil {
		return err
	}
	fmt.Printf("NUM: %d, ERR: %v", num, err)

	for _, item := range maps {
		fmt.Println(item["id"], ":", item["name"])
	}

	return nil
}

// TransOperation 事务操作
func TransOperation() error {

	tx, err := o.Begin()
	if err != nil {
		return err
	}

	user := User{Name: "slene"}
	_, err = tx.Insert(&user)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	return nil
}
