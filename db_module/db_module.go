package db_module

import (
	"github.com/jinzhu/gorm"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

const luaOrmDbName = "db_metatable"

type OrmDB struct {
	DbName string
	Tag    string
	Db     *gorm.DB
}

func RegisterOrmDbType(L *lua.LState) {
	Setup()
	mt := L.NewTypeMetatable(luaOrmDbName)
	L.SetGlobal("db_module", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newDb))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), gormMethods))
}

var gormMethods = map[string]lua.LGFunction{
	"Tag":    Tag,
	"DbName": DbName,
	"Raw":    raw,
	"Exec":   exec,
}

// Constructor
func newDb(L *lua.LState) int {
	var DbName = L.CheckString(1)
	ormDb := &OrmDB{
		DbName: DbName,
		Tag:    "初始化",
		Db:     GetDB().Table(DbName),
	}
	ud := L.NewUserData()
	ud.Value = ormDb
	L.SetMetatable(ud, L.GetTypeMetatable(luaOrmDbName))
	L.Push(ud)
	return 1
}

func raw(L *lua.LState) int {
	ormDb := checkDb(L)
	var execSql = L.CheckString(2)
	rows, err := ormDb.Db.Raw(execSql).Rows()
	ormDb.Tag = "raw"
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
	for k, _ := range vals {
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
	ormDb := checkDb(L)
	var execSql = L.CheckString(2)
	ormDb.Tag = "exec"
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

// Getter and setter for the Person#Name
func Tag(L *lua.LState) int {
	p := checkDb(L)
	L.Push(lua.LString(p.Tag))
	return 1
}

// Getter and setter for the Person#Name
func DbName(L *lua.LState) int {
	p := checkDb(L)
	L.Push(lua.LString(p.DbName))
	return 1
}

func checkDb(L *lua.LState) *OrmDB {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*OrmDB); ok {
		return v
	}
	L.ArgError(1, "person expected")
	return nil
}

func transLuaValue2Map(value lua.LValue) interface{} {
	if value.Type() == lua.LTTable {
		var deMap = make(map[string]interface{})
		var list []interface{}
		var table = value.(*lua.LTable)
		table.ForEach(func(key lua.LValue, value lua.LValue) {
			if key.Type() == lua.LTNumber {
				list = append(list, transLuaValue2Map(value))
			} else {
				deMap[key.String()] = transLuaValue2Map(value)
			}
		})
		if len(deMap) > 0 && len(list) > 0 {
			return map[string]interface{}{
				"map":  deMap,
				"list": list,
			}
		}
		if len(deMap) > 0 {
			return deMap
		}
		if len(list) > 0 {
			return list
		}
		return deMap
	} else if value.Type() == lua.LTUserData {
		var table = value.(*lua.LUserData)
		return table.Value
	} else {
		return value
	}
}
