package steam_helper

import "time"

//
var exchangeRate = map[string]float64{
	"USD ($)":    1,  // USD (доллар США)
	"GBP (£)":    2,  // GBP (фунт стерлингов)
	"EUR (€)":    3,  // EUR (евро)
	"CHF (CHF)":  4,  // CHF (швейцарский франк)
	"RUB (pуб.)": 5,  // RUB (российский рубль)
	"PLN (zł)":   6,  // PLN (польский злотый)
	"BRL (R$)":   7,  // BRL (бразильский реал)
	"JPY (¥)":    8,  // JPY (японская иена)
	"NOK (kr)":   9,  // NOK (норвежская крона)
	"IDR (Rp)":   10, // IDR (индонезийская рупия)
	"MYR (RM)":   11, // MYR (малайзийский ринггит)
	"PHP (₱)":    12, // PHP (филиппинское песо)
	"SGD (S$)":   13, // SGD (сингапурский доллар)
	"THB (฿)":    14, // THB (тайский бат)
	"VND (₫)":    15, // VND (вьетнамский донг)
	"KRW (₩)":    16, // KRW (южнокорейская вона)
	"TRY (₺)":    17, // TRY (турецкая лира)
	"UAH (₴)":    18, // UAH (украинская гривна)
	"MXN (Mex$)": 19, // MXN (мексиканское песо)
	"CAD (CDN$)": 20, // CAD (канадский доллар)
	"AUD (A$)":   21, // AUD (австралийский доллар)
	"NZD (NZ$)":  22, // NZD (новозеландский доллар)
	"CNY (元)":    23, // CNY (китайский юань)
	"INR (₹)":    24, // INR (индийская рупия)
	"CLP (CLP$)": 25, // CLP (чилийское песо)
	"PEN (S/)":   26, // PEN (перуанский соль)
	"COP (COL$)": 27, // COP (колумбийское песо)
	"ZAR (R)":    28, // ZAR (южноафриканский рэнд)
	"HKD (HK$)":  29, // HKD (гонконгский доллар)
	"TWD (NT$)":  30, // TWD (тайваньский доллар)
	"SAR (SAR)":  31, // SAR (саудовский риял)
	"AED (AED)":  32, // AED (дирхам ОАЭ)
	"SEK (kr)":   33, // SEK (шведская крона)
	"ARS (ARS$)": 34, // ARS (аргентинское песо)
	"ILS (₪)":    35, // ILS (новый израильский шекель)
	"BYN (Br)":   36, // BYN (беларуский рубль)
	"KZT (₸)":    37, // KZT (казахстанский тенге)
	"KWD (K.D.)": 38, // KWD (кувейтский динар)
	"QAR (QAR)":  39, // QAR (катарский риал)
	"CRC (₡)":    40, // CRC (коста-риканский колон)
	"UYU ($U)":   41, // UYU (уругвайское песо)
	"BGN (лв)":   42, // BGN (болгарский лев)
	"HRK (kn)":   43, // HRK (хорватская куна)
	"CZK (Kč)":   44, // CZK (чешская крона)
	"DKK (kr)":   45, // DKK (датская крона)
	"HUF (Ft)":   46, // HUF (венгерский форинт)
	"RON (L)":    47, // RON (румынский лей)
}

func TimeParse(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", t)
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

var ERROR_NO_SUCH_ELEMENT_IN_FRAME = "stale element reference: stale element reference: stale element not found in the current frame"