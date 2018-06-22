package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Набор правил для проверки корректности строки, получаемой из файлов.
var config = struct{ date, geo, zone, impressions, revenue *regexp.Regexp }{
	regexp.MustCompile(`^[-]?[0-9]{4}-(((0[13578]|(10|12))-(0[1-9]|[1-2][0-9]|3[0-1]))|(02-(0[1-9]|[1-2][0-9]))|((0[469]|11)-(0[1-9]|[1-2][0-9]|30)))$`),
	regexp.MustCompile(`[A-Z]{2}`),
	regexp.MustCompile(`^[+-]?([0-9]*[.])?[0-9]+$`),
	regexp.MustCompile(`^[0-9]*$`),
	regexp.MustCompile(`[+-]?([0-9]*[.])?[0-9]+`),
}

// checkLine проверяет полученную запись (строку из csv файла) на корректность.
func checkLine(record []string) bool {
	return config.date.MatchString(record[0]) &&
		config.geo.MatchString(record[1]) &&
		config.zone.MatchString(record[2]) &&
		config.impressions.MatchString(record[3]) &&
		config.revenue.MatchString(record[4])
}

// getFies получает список файлов из указанной директории
func getFies(searchDir string) (fileList []string, err error) {
	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			if filepath.Ext(path) == ".csv" {
				fileList = append(fileList, path)
			}
		}
		return nil
	})
	return
}

// getFileContent считывает данные со входящиего файла,
// осуществляет проверку и записывает в результирующий бинарный файл
func getFileContent(fileName string) error {
	fileErrors, err := os.Create("errors.csv")
	if err != nil {
		return err
	}
	defer fileErrors.Close()

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writerErrors := csv.NewWriter(fileErrors)
	defer writerErrors.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if !checkLine(record) {
			log.Println("Найдена некорректная запись:", record)
			err = writerErrors.Write(record)
			if err != nil {
				return err
			}
			continue
		}

		node := tree.CreateNode([]byte(strings.Join(record, ",")))
		tree.Add(node)
	}
	return nil
}
