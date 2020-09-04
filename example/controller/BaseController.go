package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

func formatSuccess(data interface{}) gin.H {
	return gin.H{
		"code": 200,
		"msg":  "ok",
		"data": data,
	}
}

func formatError(data interface{}) gin.H {
	return gin.H{
		"code": -1,
		"msg":  data,
		"data": "",
	}
}

func getNowTime(keepTime int) string {
	keepTime = keepTime / 1000
	var day = keepTime / 60 / 60 / 24
	var hour = keepTime / 60 / 60 % 24
	var minute = keepTime / 60 % 60
	var second = keepTime % 60
	return fmt.Sprint(day, "天", hour, "时", minute, "分", second, "秒")
}

func nowTime() int64 {
	return time.Now().Unix()
}

/**
 * 包含判断
 */
func Contains(a []string, word string) bool {
	for _, key := range a {
		if key == word {
			return true
		}
	}
	return false
}

/**
 * 手机号验证
 */
func isPhone(phone string) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

/**
 * 手机号加*
 */
func MobileReplaceRepl(str string) string {
	re, _ := regexp.Compile("(\\d{3})(\\d{4})(\\d{4})")
	return re.ReplaceAllString(str, "$1****$3")
}
