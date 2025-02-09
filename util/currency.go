package util

// Константы для всех поддерживаемых валю
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	RUB = "RUB"
)

// IsSupportedCurrency возвращает true, если валюта поддерживается
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, RUB:
		return true
	}
	return false
}