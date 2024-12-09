package phone

import (
	"fmt"
	"strconv"
)

func Present(phone int64) string {
	phoneStr := strconv.FormatInt(phone, 10)

	if len(phoneStr) != 11 {
		return "Неверный номер телефона"
	}

	formatted := fmt.Sprintf("+7 (%s) %s-%s-%s",
		phoneStr[1:4],  // 915
		phoneStr[4:7],  // 159
		phoneStr[7:9],  // 67
		phoneStr[9:11], // 81
	)

	return formatted
}
