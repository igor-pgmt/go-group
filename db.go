package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	mysql "github.com/go-sql-driver/mysql"
)

// convertToStr конвертирует результирующий файл в csv
func convertToStr() error {
	file, err := os.Create("test.csv")
	defer file.Close()
	if err != nil {
		return err
	}

	for i := int64(0); ; i += 64 {
		m := new(Node)
		data, err := readBytes(tree.File, i)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, m)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		a := string(bytes.Trim([]byte(m.Data[:]), "\x00")) + "\n"
		file.WriteString(a)
	}
	return nil
}

// importCSV Импортирует csv файл в базу даных в таблицу t2
func importCSV() error {
	filePath := "./test.csv"
	mysql.RegisterLocalFile(filePath)
	_, err := db.Exec("LOAD DATA LOCAL INFILE '" + filePath + "' INTO TABLE t2 FIELDS TERMINATED BY ','")
	return err
}

// mergeTables копирует данные из таблицы t2 в t1
func mergeTables() error {
	_, err := db.Exec("insert t1 select * from t2 on duplicate key UPDATE t1.impressions = t1.impressions + values(impressions), t1.revenue = t1.revenue + values(revenue);")
	return err
}
