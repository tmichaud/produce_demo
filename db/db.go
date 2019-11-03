package db

import (
	"produce_demo/common"
	"strings"
	"sync"
)

var mutex = &sync.Mutex{}

var rows = map[string]common.Produce{
	"A12T-4GH7-QPL9-3N4M": common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
	"E5T6-9UI3-TH15-QR88": common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
	"YRT6-72AS-K736-L4AR": common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
	"TQ4C-VV6T-75ZX-1RMR": common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
}

// Concurrent Add of Produce
func Add(p common.Produce, outputChannel chan<- common.Result) {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.ToUpper(p.ProduceCode)
	_, ok := rows[key]
	if ok {
		// key exists
		outputChannel <- common.Result{Prod: p, Err: key + " already exists", Count: 0}
	} else {
		rows[key] = p
		outputChannel <- common.Result{Prod: rows[key], Err: "", Count: 1}
	}
}

// Concurrent Delete
func Delete(produceCode string, outputChannel chan<- common.Result) {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.ToUpper(produceCode)
	_, ok := rows[key]
	if ok {
		outputChannel <- common.Result{Prod: rows[key], Err: "", Count: 1}
		delete(rows, key)
	} else {
		outputChannel <- common.Result{Prod: common.Produce{}, Err: "Row not found", Count: 0}
	}
}

// Concurrent DeleteRow
func DeleteRow(row common.Produce, outputChannel chan<- common.Result) {
	Delete(row.ProduceCode, outputChannel)
}

// Concurrent Fetch
func Fetch(outputChannel chan<- common.Result) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(rows) == 0 {
		outputChannel <- common.Result{Prod: common.Produce{}, Err: "Row not found", Count: 0}
	} else {
		for k := range rows {
			outputChannel <- common.Result{Prod: rows[k], Err: "", Count: 1}
		}
	}
	close(outputChannel)
}

// Concurrent FetchByProduceCode
func FetchByProduceCode(produceCode string, outputChannel chan<- common.Result) {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.ToUpper(produceCode)
	prod, ok := rows[key]
	if ok {
		outputChannel <- common.Result{Prod: prod, Err: "", Count: 1}
	} else {
		outputChannel <- common.Result{Prod: common.Produce{}, Err: "Row not found", Count: 0}
	}
	close(outputChannel)
}
