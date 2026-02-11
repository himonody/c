package utils

import (
	"fmt"
	"time"
)

const MYDateLayout = "02/01/2006 15:04:05"

func GenerateWithdrawBizID(userId int64) string {
	now := time.Now()
	ms := now.Nanosecond() / 1e6
	return fmt.Sprintf("WD%d%s%03d", userId, now.Format("20060102150405"), ms)
}
