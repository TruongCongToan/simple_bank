package util

// constants for all supported currency formats
const (
	USD = "USD"
	VND = "VND"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportCurrency(currency string) bool {
	switch currency {
	case USD, EUR, VND, CAD:
		return true
	}
	return false
}
