# Реализация тестового задания
Задание располагается [здесь](https://gist.github.com/pavelkdev/4a37e7cb814bf757ed378f0c736d1c61).
Текст задания:
>Создать консольное приложение, которое на вход получает путь к папке. По этому пути находится неизвестное число папок с неизвестной глубиной вложенности. В папках находятся файлы разных типов, среди которых есть файлы с расширением csv, в них хранятся данные в следующем формате:
```
date,geo,zone,impressions,revenue
2018-01-01,RU,1111111,100,0.25
...
```
>Необходимо проитерировать все csv файлы и, собрать данные сгруппировав по полям date,geo,zone. Часть строк или файлов целиком могут иметь некорректный формат данных. В таком случае нужно вывести сообщение об ошибке в консоль и проигнорировать строку или файл с некорректными данными.
>Финальный, сгруппированный набор данных может содержать сотни тысяч строк. Необходимо записать эти данные в БД. В БД уже могут быть записи с таким же ключом (date,geo,zone). В таком случае необходимо обновить строки просуммировав значения impressions и revenue хранящиеся в бд со значениями, полученными в ходе агрегации.
___
### Допущения реализации
Исходя из формулировки:
>сгруппированный набор данных может содержать сотни тысяч строк

предполагается, что даты (параметр **date**) могут могут быть в широком диапазоне, включая "до нашей эры". Даты подобного вида могут иметь различный формат представления (префиксы или постфиксы BC, AD и т.п.); мы будем использовать соответствующую дату с префиксом "-".

Так же, предполагается, что: 
**geo** - параметр всегда имеет две буквы верхнего регистра определённого множества корректных комбинаций;
**zone** - положительное целое число, максимальный объём - 64 бит;
**impressions** - целое число из диапазона [0:100];
**revenue** - число с плавающей точкой c двумя знаками после запятой из диапазона [-1:1];

### Предварительные требования
Программа предполагает наличие базы данных MySQL на localhost с двумя заранее созданными таблицами: t1 и t2:
```
CREATE TABLE `test`.`t1` ( `date` VARCHAR(11) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL , `geo` VARCHAR(2) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL , `zone` INT NOT NULL , `impressions` INT NOT NULL , `revenue` FLOAT NOT NULL, unique key(`date`,`geo`,`zone`)) ENGINE = InnoDB CHARSET=utf8 COLLATE utf8_general_ci;
```
```
CREATE TABLE `test`.`t2` ( `date` VARCHAR(11) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL , `geo` VARCHAR(2) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL , `zone` INT NOT NULL , `impressions` INT NOT NULL , `revenue` FLOAT NOT NULL, unique key(`date`,`geo`,`zone`)) ENGINE = InnoDB CHARSET=utf8 COLLATE utf8_general_ci;
```
Таблица t1 используется для хранения результирующих данных. Она может быть пустой или уже содержать некоторые данные (заранее корректного формата). Пример данных [здесь](https://github.com/igor-pgmt/go-group/blob/master/t1.sql); по желанию можно произвести импорт в таблицу t1 перед началом тестирования программы.

### Использование программы
Скачать и скомпилировать:
```terminal
git clone https://github.com/igor-pgmt/go-group.git
cd go-group
go build
```
После этого программа готова к работе.

Программа обладает возможностью генерировать тестовые данные (реализовано в файле [datagen.go](https://github.com/igor-pgmt/go-group/blob/master/datagen.go)). С помощью параметра **tdata**, указав директорию, программа сгенерирует тестовые данные в этой директории (создаст директорию, если отсутствует):
```terminal
./go-group -tdata test
```

Поиск файлов csv, группировка данных и импорт в БД осуществляется с помощью одновременного указания ключей **dir** (директория для поиска файлов csv) и **cred** (учётные данные для доступа к БД в формате login:pass@/dbname):
```terminal
./go-group -dir test -cred root:123@/test
```
Следует указать БД, содержащую таблицы, t1 и t2, созданные ранее.

При запуске данной команды, программа произведёт следующие действия:
- Поиск файлов csv в указанной директории
- Построчное считывание данных из файлов, проверка данных на корректность и запись в результирующий файл test.bin с учётом логики группировки данных, указанной в задании; при этом, некорректные данные будут отображены в терминале и выведены в файл errors.csv
- Конвертирование файла test.bin в файл test.csv
- Импорт файла test.csv в таблицу t2 указанной пользователем БД (посредством [LOAD DATA INFILE](https://dev.mysql.com/doc/refman/8.0/en/load-data.html))
- Перенос данных из таблицы t2 в таблицу t1 с учётом условий группировки данных

Результатом работы является содержимое таблицы t1.
После работы программы остаются доступными файлы test.bin, test.csv, errors.csv.

### Возможные доработки/изменения для обработки реальных данных
- Изменение типов данных и\или формата представления данных, получаемых из исходных csv
- В зависимости от того или иного объёма данных или операций сложения данных, возможно, потребуется увеличить объём используемых типов данных; например, вместо int64 использовать bigInt и т.п.