package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Дерево
type Tree struct {
	Root    *Node
	File    *os.File
	Counter int64
}

// Узел дерева
type Node struct {
	Data  [40]byte
	This  int64
	Left  int64
	Right int64
}

// CreateNode Создаёт узел дерева для дальнейшей вставки в дерево.
func (t *Tree) CreateNode(a0 []byte, a ...int64) *Node {
	node := new(Node)
	copy(node.Data[:], a0)
	if len(a) > 0 {
		node.This = a[0]
	}
	if len(a) > 1 {
		node.Left = a[1]
	}
	if len(a) > 2 {
		node.Right = a[2]
	}
	return node
}

// getNode возвращает узел по указанному идентификатору.
// Для каждого узла идентификатором является номер строки с данными.
func (t *Tree) getNode(from int64) *Node {
	node := new(Node)
	data, err := readBytes(t.File, from)
	if err == io.EOF {
		return nil
	} else if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, node)
	if err != nil {
		panic(err)
	}
	return node
}

// getRootNode возвращает корневой узел дерева.
func (t *Tree) getRootNode() *Node {
	return t.getNode(0)
}

// Find ищет указанный узел дерева и возвращает найденный узел и его родителя.
// В случае, если узел отсутствует, возвращается <nil> и возможный родительский узел.
func (t *Tree) Find(d []byte) (*Node, *Node, int64) {
	parentNode, currentNode := t.Root, t.Root
	var counter int64
	a := strings.Split(string(d), ",")
	aet := strings.Join(a[:3], ",")

	for currentNode != nil {
		counter++
		if aet > strings.Join(strings.Split(string(currentNode.Data[:]), ",")[:3], ",") {
			if currentNode.Right == 0 {
				return nil, currentNode, counter
			}
			parentNode = currentNode
			currentNode = t.getNode(currentNode.Right)

		} else if aet < strings.Join(strings.Split(string(currentNode.Data[:]), ",")[:3], ",") {
			if currentNode.Left == 0 {
				return nil, currentNode, counter
			}
			parentNode = currentNode
			currentNode = t.getNode(currentNode.Left)

		} else {
			return currentNode, parentNode, counter
		}
	}
	return nil, parentNode, counter
}

// Add добавляет узел в дерево (записывает в результирующий файл)
// и возвращает идентификатор добавленного узла
func (t *Tree) Add(n *Node) int64 {
	findedNode, parentNode, _ := t.Find(n.Data[:])
	if findedNode != nil {
		fn := strings.Split(string(findedNode.Data[:]), ",")
		nd := strings.Split(string(n.Data[:]), ",")
		aa1, _ := strconv.Atoi(fn[3])
		aa2, _ := strconv.Atoi(nd[3])
		fn[3] = strconv.Itoa(aa1 + aa2)
		aa11, _ := strconv.ParseFloat(string(bytes.Trim([]byte(fn[4]), "\x00")), 32)
		aa22, _ := strconv.ParseFloat(string(bytes.Trim([]byte(nd[4]), "\x00")), 32)
		fn[4] = strconv.FormatFloat(aa11+aa22, 'f', -1, 32)
		strj := strings.Join(fn, ",")
		copy(findedNode.Data[:], strj)
		writeBytes(t.File, findedNode)
		return findedNode.This
	} else if findedNode == nil && parentNode != nil {
		newNode := new(Node)
		newNode.Data = n.Data
		newNode.This = t.Counter
		newNode.Left = n.Left
		newNode.Right = n.Right
		if strings.Join(strings.Split(string(n.Data[:]), ",")[:3], ",") > strings.Join(strings.Split(string(parentNode.Data[:]), ",")[:3], ",") {
			parentNode.Right = newNode.This
		} else {
			parentNode.Left = newNode.This
		}
		t.Counter += 64
		writeBytes(t.File, newNode)
		writeBytes(t.File, parentNode)

		return newNode.This
	}
	newNode := new(Node)
	newNode.Data = n.Data
	newNode.This = t.Counter
	newNode.Left = n.Left
	newNode.Right = n.Right
	t.Root = newNode
	t.Counter += 64
	writeBytes(t.File, newNode)

	return newNode.This
}

// writeBytes записывает данные в файл
func writeBytes(file *os.File, n *Node) error {
	buf := bytes.NewBuffer(make([]byte, 0, 64))
	binary.Write(buf, binary.BigEndian, n)
	_, err := file.WriteAt(buf.Bytes(), n.This)
	if err != nil {
		return err
	}
	return nil
}

// readBytes считывает данные из файла
func readBytes(file *os.File, from int64) ([]byte, error) {
	bytes := make([]byte, 64)
	_, err := file.ReadAt(bytes, from)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// PrintTree - отладочная функция для просмотра содержимого дерева
func (t *Tree) PrintTree() {
	for i := int64(0); ; i += 64 {
		m := new(Node)
		data, err := readBytes(t.File, i)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, m)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Printf("%v, | %d, %d, %d\n", string(m.Data[:]), m.This, m.Left, m.Right)
	}
}
