package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"os"
	"regexp"
)

// Параметры командной строки

var tdata string // Директория для генерации тестовых данных
var dir string   // Директория поиска файлов csv
var cred string  // Учётные данные для БД

// Дерево данных
var tree Tree

// База данных
var db *sql.DB

func main() {

	// Проверка аргументов программы
	err := checkFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Генерация тестовых данных.
	if tdata != "" {
		log.Println("Генерация тестовых данных...")

		err := GenerateData(5, 5)
		if err != nil {
			panic(err)
		}
		log.Println("Работа завершена")
		return
	}

	// Получение списка файлов указанной директории.
	log.Println("Директория для поиска файлов:", dir)
	fileList, err := getFies(dir)
	if err != nil {
		log.Fatal("Нет доступа к директории или директория отсутствует.")
	}
	if len(fileList) == 0 {
		log.Fatal("Файлы csv не найдены.")
	}

	// Cоединение с базой данных.
	db, err = sql.Open("mysql", cred)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверка соединния
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Создание результирующего файла.
	tree.File, err = os.Create("test.bin")
	defer tree.File.Close()
	if err != nil {
		panic(err)
	}

	// Обработка файлов.
	for _, fileName := range fileList {
		log.Println("Обрабатывается файл:", fileName)
		getFileContent(fileName)
	}

	// Конвертирование результирующего файла в csv.
	err = convertToStr()
	if err != nil {
		panic(err)
	}

	// Импорт результирующего csv в БД.
	err = importCSV()
	log.Println("Импорт в базу данных...")
	if err != nil {
		panic(err)
	}

	// Копирование импортированных данных в общую таблицу.
	err = mergeTables()
	log.Println("Импорт в общую таблицу...")
	if err != nil {
		panic(err)
	}

	log.Println("Работа завершена")
}

// checkFlags регистрирует и проверяет аргументы приложения
func checkFlags() error {
	flag.StringVar(&dir, "dir", "", "[string] Указать директорию для поиска файлов.")
	flag.StringVar(&tdata, "tdata", "", "[string] Сгенерировать тестовые данные в указанную папку.")
	flag.StringVar(&cred, "cred", "", "[string] Указать данные доступа к БД в формате login:pass@/dbname")
	flag.Parse()
	if tdata == "" && cred == "" {
		flag.Usage()
		return errors.New("Обязательно должен быть указан либо параметр cred, либо tdata и dir.")
	}
	if tdata != "" && (cred != "" || dir != "") {
		flag.Usage()
		return errors.New("Возможно одновременное использование либо параметра tdata, либо cred и dir.")
	}
	if tdata == "" && (cred == "" || dir == "") {
		flag.Usage()
		return errors.New("Необходимо одновременное использование параметров cred и dir.")
	}

	if cred != "" {
		r := regexp.MustCompile(`.*\:.*\@\/.*`)
		if !r.MatchString(cred) {
			flag.Usage()
			return errors.New("Формат доступа к БД: login:pass@/dbname")
		}
	}

	if len(dir) > 0 && dir[len(dir)-1:] != "/" {
		dir += "/"
	}
	if len(tdata) > 0 && tdata[len(tdata)-1:] != "/" {
		tdata += "/"
	}
	return nil
}
