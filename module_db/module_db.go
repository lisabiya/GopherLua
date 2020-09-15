package module_db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"time"
)

//数据库模块
type ModuleDb struct {
	DbCreateCallBack func(*gorm.DB)
	OrmDBs           []*OrmDB
}

type OrmDB struct {
	DbPath string
	Tag    string
	Db     *gorm.DB
}

const metatableName = "db_metatable"

func (db ModuleDb) RegisterType(L *lua.LState) {
	mt := L.NewTypeMetatable(metatableName)
	L.SetGlobal("db_module", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newObject(&db)))
	L.SetField(mt, "closeDbByTag", L.NewFunction(closeDbByTag(&db)))
	L.SetField(mt, "closeAllDb", L.NewFunction(closeAllDb(&db)))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), gormMethods))
}

func (db ModuleDb) Close() {
	for _, orm := range db.OrmDBs {
		if orm.Db != nil {
			err := orm.Db.Close()
			if err != nil {
				fmt.Printf("关闭指定数据库%s成功,错误信息%s", orm.Tag, err.Error())
			} else {
				fmt.Printf("关闭指定数据库%s成功", orm.Tag)
			}
			orm.Db = nil
		}
	}
	db.OrmDBs = nil
}

func newObject(moduleDb *ModuleDb) lua.LGFunction {
	return func(state *lua.LState) int {
		var dialect = state.CheckString(1)
		var path = state.CheckString(2)
		var enableLog = state.CheckBool(3)
		var tag = state.CheckString(4)

		var db = setup(dialect, path, enableLog)

		var ormDb = &OrmDB{
			DbPath: path,
			Tag:    tag,
			Db:     db,
		}
		moduleDb.OrmDBs = append(moduleDb.OrmDBs, ormDb)
		fmt.Println("数据库数量", len(moduleDb.OrmDBs))
		if moduleDb.DbCreateCallBack != nil {
			moduleDb.DbCreateCallBack(db)
		}
		ud := state.NewUserData()
		ud.Value = ormDb
		state.SetMetatable(ud, state.GetTypeMetatable(metatableName))
		state.Push(ud)
		return 1
	}
}

func closeDbByTag(moduleDb *ModuleDb) lua.LGFunction {
	return func(state *lua.LState) int {
		var tag = state.CheckString(1)
		for _, orm := range moduleDb.OrmDBs {
			if orm.Db != nil && orm.Tag == tag {
				err := orm.Db.Close()
				orm.Db = nil
				if err != nil {
					fmt.Printf("关闭指定数据库%s成功,错误信息%s", orm.Tag, err.Error())
				} else {
					fmt.Printf("关闭指定数据库%s成功", orm.Tag)
				}

			}
		}
		var newArr []*OrmDB
		for _, orm := range moduleDb.OrmDBs {
			if orm.Db != nil {
				newArr = append(newArr, orm)
			}
		}
		moduleDb.OrmDBs = newArr
		return 1
	}
}

func closeAllDb(moduleDb *ModuleDb) lua.LGFunction {
	return func(state *lua.LState) int {
		moduleDb.Close()
		state.Push(lua.LString("清空完毕"))
		return 1
	}
}

func setup(dialect, path string, enableLog bool) *gorm.DB {
	db, err := gorm.Open(dialect, path)
	if err != nil {
		panic(err.Error())
	}
	db.DB().SetConnMaxLifetime(80 * time.Second) // 设置链接重置时间
	db.LogMode(enableLog)
	return db
}

//**********元表方法***********
var gormMethods = map[string]lua.LGFunction{
	"Tag":     tag,
	"Raw":     raw,
	"Exec":    exec,
	"CloseDB": closeDB,
}

func raw(L *lua.LState) int {
	ormDb := checkObject(L)
	var execSql = L.CheckString(2)
	rows, err := ormDb.Db.Raw(execSql).Rows()
	if err != nil {
		L.Push(lua.LNumber(1))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k := range vals {
		scans[k] = &vals[k]
	}
	var table = lua.LTable{}
	i := 1
	for rows.Next() {
		//rows
		//填充数据
		_ = rows.Scan(scans...)
		//每行数据
		row := make(map[string]interface{})
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		table.Insert(i, luar.New(L, row))
		i++
	}
	L.Push(lua.LNumber(0))
	L.Push(&table)
	return 2
}

func exec(L *lua.LState) int {
	ormDb := checkObject(L)
	var execSql = L.CheckString(2)
	err := ormDb.Db.Exec(execSql).Error
	if err != nil {
		L.Push(lua.LNumber(1))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(0))
	L.Push(lua.LString("成功"))
	return 2
}

func tag(L *lua.LState) int {
	p := checkObject(L)
	L.Push(lua.LString(p.Tag))
	return 1
}

func closeDB(L *lua.LState) int {
	p := checkObject(L)
	err := p.Db.Close()
	if err != nil {
		L.Push(lua.LNumber(1))
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString("关闭成功"))
	}
	return 2
}

func checkObject(L *lua.LState) *OrmDB {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*OrmDB); ok {
		return v
	}
	L.ArgError(1, "person expected")
	return nil
}
