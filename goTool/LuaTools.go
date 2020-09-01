package goTool

import lua "github.com/yuin/gopher-lua"

func TransLuaValue2Map(value lua.LValue) interface{} {
	switch value.Type() {
	case lua.LTNumber:
		return value.(lua.LNumber)
	case lua.LTString:
		return string(value.(lua.LString))
	case lua.LTBool:
		return bool(value.(lua.LBool))
	case lua.LTTable:
		var deMap = make(map[string]interface{})
		var list []interface{}
		var table = value.(*lua.LTable)
		table.ForEach(func(key lua.LValue, value lua.LValue) {
			if key.Type() == lua.LTNumber {
				list = append(list, TransLuaValue2Map(value))
			} else {
				deMap[key.String()] = TransLuaValue2Map(value)
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
	case lua.LTUserData:
		var table = value.(*lua.LUserData)
		return table.Value
	default:
		return value
	}
}
