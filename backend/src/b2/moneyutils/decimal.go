package moneyutils

import (
	"b2/errors"
	"fmt"
	"strconv"
)

var ccyDefs = map[string]uint64{
	"AED": 100,
	"AFN": 100,
	"ALL": 100,
	"AMD": 100,
	"ANG": 100,
	"AOA": 100,
	"ARS": 100,
	"AUD": 100,
	"AWG": 100,
	"AZN": 100,
	"BAM": 100,
	"BBD": 100,
	"BDT": 100,
	"BGN": 100,
	"BHD": 1000,
	"BIF": 100,
	"BMD": 100,
	"BND": 100,
	"BOB": 100,
	"BRL": 100,
	"BSD": 100,
	"BTN": 100,
	"BWP": 100,
	"BYN": 100,
	"BZD": 100,
	"CAD": 100,
	"CDF": 100,
	"CHF": 100,
	"CLP": 100,
	"CNY": 10,
	"COP": 100,
	"CRC": 100,
	"CUC": 100,
	"CUP": 100,
	"CVE": 100,
	"CZK": 100,
	"DJF": 100,
	"DKK": 100,
	"DOP": 100,
	"DZD": 100,
	"EGP": 100,
	"ERN": 100,
	"ETB": 100,
	"EUR": 100,
	"FJD": 100,
	"FKP": 100,
	"GBP": 100,
	"GEL": 100,
	"GHS": 100,
	"GIP": 100,
	"GMD": 100,
	"GNF": 100,
	"GTQ": 100,
	"GYD": 100,
	"HKD": 100,
	"HNL": 100,
	"HRK": 100,
	"HTG": 100,
	"HUF": 100,
	"IDR": 100,
	"ILS": 100,
	"INR": 100,
	"IQD": 1000,
	"IRR": 100,
	"ISK": 100,
	"JMD": 100,
	"JOD": 100,
	"JPY": 100,
	"KES": 100,
	"KGS": 100,
	"KHR": 100,
	"KMF": 100,
	"KPW": 100,
	"KRW": 100,
	"KWD": 1000,
	"KYD": 100,
	"KZT": 100,
	"LAK": 100,
	"LBP": 100,
	"LKR": 100,
	"LRD": 100,
	"LSL": 100,
	"LYD": 1000,
	"MAD": 100,
	"MDL": 100,
	"MGA": 5,
	"MKD": 100,
	"MMK": 100,
	"MNT": 100,
	"MOP": 100,
	"MRU": 5,
	"MUR": 100,
	"MVR": 100,
	"MWK": 100,
	"MXN": 100,
	"MYR": 100,
	"MZN": 100,
	"NAD": 100,
	"NGN": 100,
	"NIO": 100,
	"NOK": 100,
	"NPR": 100,
	"NZD": 100,
	"OMR": 1000,
	"PAB": 100,
	"PEN": 100,
	"PGK": 100,
	"PHP": 100,
	"PKR": 100,
	"PLN": 100,
	"PYG": 100,
	"QAR": 100,
	"RON": 100,
	"RSD": 100,
	"RUB": 100,
	"RWF": 100,
	"SAR": 100,
	"SBD": 100,
	"SCR": 100,
	"SDG": 100,
	"SEK": 100,
	"SGD": 100,
	"SHP": 100,
	"SLL": 100,
	"SOS": 100,
	"SRD": 100,
	"SSP": 100,
	"STN": 100,
	"SYP": 100,
	"SZL": 100,
	"THB": 100,
	"TJS": 100,
	"TMT": 100,
	"TND": 1000,
	"TOP": 100,
	"TRY": 100,
	"TTD": 100,
	"TWD": 100,
	"TZS": 100,
	"UAH": 100,
	"UGX": 100,
	"USD": 100,
	"UYU": 100,
	"UZS": 100,
	"VES": 100,
	"VND": 10,
	"VUV": 0,
	"WST": 100,
	"XAF": 100,
	"XCD": 100,
	"XOF": 100,
	"XPF": 100,
	"YER": 100,
	"ZAR": 100,
	"ZMW": 100,
}

// CurrencyAmount returns a float64 that represents the currency in its higher base
// this might contain errors inherent in float math
func CurrencyAmount(amount int64, ccy string) (float64, error) {
	multiple, ok := ccyDefs[ccy]
	if !ok {
		return 0, errors.New("CCY definition not found for "+ccy, nil, "moneyutils.CurrencyAmount")
	}
	if multiple == 0 {
		return float64(amount), nil
	}
	return float64(amount) / float64(multiple), nil
}

// CurrencyAmountPrint returns a string representation of the amount in the given currency
// with the amount rounded to the nearest minor unit
func CurrencyAmountPrint(amount int64, ccy string) (string, error) {
	multiple, ok := ccyDefs[ccy]
	if !ok {
		return "", errors.New("CCY definition not found for "+ccy, nil, "moneyutils.CurrencyAmount")
	}
	switch multiple {
	case 0:
		return fmt.Sprintf("%d", amount), nil
	case 5:
		return fmt.Sprintf("%.1f", float64(amount)/float64(multiple)), nil
	case 10:
		return fmt.Sprintf("%.1f", float64(amount)/float64(multiple)), nil
	case 100:
		return fmt.Sprintf("%.2f", float64(amount)/float64(multiple)), nil
	case 1000:
		return fmt.Sprintf("%.3f", float64(amount)/float64(multiple)), nil
	default:
		panic("Someone updated the ccy definitions and didn't update CurrencyAmountPrint")
	}
}

// CurrencyAmountPrint returns a string representation of the absolute amount in the given currency
// with the amount rounded to the nearest minor unit
func AbsCurrencyAmountPrint(amount int64, ccy string) (string, error) {
	multiple, ok := ccyDefs[ccy]
	if !ok {
		return "", errors.New("CCY definition not found for "+ccy, nil, "moneyutils.CurrencyAmount")
	}
	if amount < 0 {
		amount *= -1
	}
	switch multiple {
	case 0:
		return fmt.Sprintf("%d", amount), nil
	case 5:
		return fmt.Sprintf("%.1f", float64(amount)/float64(multiple)), nil
	case 10:
		return fmt.Sprintf("%.1f", float64(amount)/float64(multiple)), nil
	case 100:
		return fmt.Sprintf("%.2f", float64(amount)/float64(multiple)), nil
	case 1000:
		return fmt.Sprintf("%.3f", float64(amount)/float64(multiple)), nil
	default:
		panic("Someone updated the ccy definitions and didn't update CurrencyAmountPrint")
	}
}

// CurrencyFromString returns an int representation of a decimal formatted currency
func CurrencyFromString(amount, ccy string) (int64, error) {
	multiple, ok := ccyDefs[ccy]
	if !ok {
		return 0, errors.New("CCY definition not found for "+ccy, nil, "moneyutils.CurrencyFromString")
	}
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, errors.Wrap(err, "moneyutils.CurrencyFromString")
	}
	if multiple == 0 {
		return int64(val), nil
	}
	return int64(val * float64(multiple)), nil
}
