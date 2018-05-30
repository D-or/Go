/*
 * Revision History:
 *     Initial: 2018/05/25      Lin Hao
 */

package utils

import (
	"fmt"
	"time"
)

// Now get now time formatted to "yaer_month_day_hour_minute_second"
func Now() string {
	now := time.Now()

	Y, M, D := now.Date()
	h := now.Hour()
	m := now.Minute()
	s := now.Second()

	timeStr := fmt.Sprintf("%d_%d_%d_%d_%d_%d", Y, M, D, h, m, s)

	return timeStr
}
