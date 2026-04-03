package luhn

import (
	"strconv"
	"strings"
)

func LuhnCheck(number string) bool {
	// Удаляем все пробелы из номера
	number = strings.ReplaceAll(number, " ", "")

	// Проверяем, что номер состоит только из цифр
	if _, err := strconv.Atoi(number); err != nil {
		return false
	}

	sum := 0
	alternate := false

	// Идем по цифрам справа налево
	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}

		sum += digit
		alternate = !alternate
	}

	// Номер валиден, если сумма делится на 10 без остатка
	return sum%10 == 0
}
