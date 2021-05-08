package user

import (
	"time"
)

const INTERVAL = int64(9000000)

type checker struct {
	time      int64
	user_id   int64
	checkType int8
	code      string
}

var storedTimer []map[int64]checker

// 再使用make函数创建一个非nil的map，nil map不能赋值

func removeTimeChecker(user_id int64, checkType int8) {
	delete(storedTimer[checkType], user_id)
}

func getCurrentTimeStamp() int64 {
	currentTime := time.Now()
	currentTimeStamp := currentTime.Unix()
	return currentTimeStamp
}

func removePastTime() {
	currentTimeStamp := getCurrentTimeStamp()
	for checkType, c := range storedTimer {
		for uid, cc := range c {
			if cc.time+INTERVAL < currentTimeStamp {
				delete(storedTimer[checkType], uid)
			}
		}
	}
}

func checkValid(user_id int64, checkType int8) bool {
	currentTimeStamp := getCurrentTimeStamp()
	givenTime := getTime(user_id, checkType)
	if givenTime+INTERVAL > currentTimeStamp {
		return false
	}
	return true
}

func getTime(user_id int64, checkType int8) int64 {
	if storedTimer == nil {
		storedTimer = make([]map[int64]checker, 3)
		for i := int(0); i < 3; i++ {
			if storedTimer[i] == nil {
				storedTimer[i] = make(map[int64]checker)
			}
		}
		return -1
	}
	value, ok := storedTimer[checkType][user_id]
	if !ok {
		return -1
	}
	return value.time
}
