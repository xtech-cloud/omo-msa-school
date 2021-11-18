package cache

func checkIDCard(number string) bool {
	// 对身份证进行简单的校验
	if len(number) != 18 {
		return false
	}
	// 系数
	coefficient := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkDigitMap := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	num := 0
	for k, v := range number[:17] {
		num += int((v - 48)) * coefficient[k]
	}
	num %= 11
	if checkDigitMap[num] != string(number[17:18]) {
		return false
	}
	return true
}
