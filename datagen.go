package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Данные для параметра geo в файле csv
var countryCodes = []string{
	"AF", "AX", "AL", "DZ", "AS", "AD", "AO", "AI", "AQ", "AG", "AR",
	"AM", "AW", "AU", "AT", "AZ", "BS", "BH", "BD", "BB", "BY", "BE",
	"BZ", "BJ", "BM", "BT", "BO", "BQ", "BA", "BW", "BV", "BR", "IO",
	"BN", "BG", "BF", "BI", "KH", "CM", "CA", "CV", "KY", "CF", "TD",
	"CL", "CN", "CX", "CC", "CO", "KM", "CG", "CD", "CK", "CR", "CI",
	"HR", "CU", "CW", "CY", "CZ", "DK", "DJ", "DM", "DO", "EC", "EG",
	"SV", "GQ", "ER", "EE", "ET", "FK", "FO", "FJ", "FI", "FR", "GF",
	"PF", "TF", "GA", "GM", "GE", "DE", "GH", "GI", "GR", "GL", "GD",
	"GP", "GU", "GT", "GG", "GN", "GW", "GY", "HT", "HM", "VA", "HN",
	"HK", "HU", "IS", "IN", "ID", "IR", "IQ", "IE", "IM", "IL", "IT",
	"JM", "JP", "JE", "JO", "KZ", "KE", "KI", "KP", "KR", "KW", "KG",
	"LA", "LV", "LB", "LS", "LR", "LY", "LI", "LT", "LU", "MO", "MK",
	"MG", "MW", "MY", "MV", "ML", "MT", "MH", "MQ", "MR", "MU", "YT",
	"MX", "FM", "MD", "MC", "MN", "ME", "MS", "MA", "MZ", "MM", "NA",
	"NR", "NP", "NL", "NC", "NZ", "NI", "NE", "NG", "NU", "NF", "MP",
	"NO", "OM", "PK", "PW", "PS", "PA", "PG", "PY", "PE", "PH", "PN",
	"PL", "PT", "PR", "QA", "RE", "RO", "RU", "RW", "BL", "SH", "KN",
	"LC", "MF", "PM", "VC", "WS", "SM", "ST", "SA", "SN", "RS", "SC",
	"SL", "SG", "SX", "SK", "SI", "SB", "SO", "ZA", "GS", "SS", "ES",
	"LK", "SD", "SR", "SJ", "SZ", "SE", "CH", "SY", "TW", "TJ", "TZ",
	"TH", "TL", "TG", "TK", "TO", "TT", "TN", "TR", "TM", "TC", "TV",
	"UG", "UA", "AE", "GB", "US", "UM", "UY", "UZ", "VU", "VE", "VN",
	"VG", "VI", "WF", "EH", "YE", "ZM", "ZW",
}

// Ошибочные данные
var countryCodesErr = []string{"11", "ww", "ZA", "1.02", "qd", "zm", "or ", "  ", "ss", "d ", " g", "fds", "dtt"}

var lcc = len(countryCodes) + len(countryCodesErr)

// generateDate Генерирует случайную дату
func generateDate(minYear, maxYear int) string {
	rand.Seed(time.Now().UnixNano())
	min := time.Date(minYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(maxYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).Local().Format("2006-01-02")
}

// generateCC генерирует случайный код страны для параметра geo
func generateCC() string {
	countryCodes = append(countryCodes, countryCodesErr...)
	rand.Seed(time.Now().UnixNano())
	return countryCodes[rand.Intn(lcc)]
}

// randomInt генерирует случайное число int из указанного диапазона; для удобства сразу возвращает строку.
func randomInt(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(max+1-min) + min)
}

// randomFloat генерирует случайное число float из указанного диапазона; для удобства сразу возвращает строку.
func randomFloat(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%.2f", rand.Float64()*float64(rand.Intn(max+1-min)+min))
}

// generateFile созадёт тестовый файл с указанными названием и количеством строк.
func generateFile(fileName string, lines int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < lines; i++ {
		value := []string{generateDate(-2000, 2000), generateCC(), randomInt(1, 1000000), randomInt(0, 100), randomFloat(-1, 1)}
		writer.Write(value)
	}

	return nil
}

// GenerateData генерирует папки с тестовыми файлами.
func GenerateData(folders, files int) error {
	err := os.MkdirAll(tdata, os.ModePerm)
	if err != nil {
		return err
	}
	for i := 0; i < folders; i++ {
		dirName := randStringRunes(8)
		err := os.MkdirAll(tdata+dirName, os.ModePerm)
		if err != nil {
			return err
		}
		for j := 0; j < files; j++ {
			fileName := randStringRunes(8) + ".csv"
			err := generateFile(tdata+dirName+"/"+fileName, 1000)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// randStringRunes генерирует случайную строку
// (используется в создании случайных имён файлов и папок).
func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
