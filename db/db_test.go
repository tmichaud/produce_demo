// db simulates a database using a single array
package db

import (
	"produce_demo/common"
	"sort"
	"strings"
	"testing"
)

func resetRows() {
	mutex.Lock()
	defer mutex.Unlock()
	rows = map[string]common.Produce{
		"A12T-4GH7-QPL9-3N4M": common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		"E5T6-9UI3-TH15-QR88": common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		"YRT6-72AS-K736-L4AR": common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		"TQ4C-VV6T-75ZX-1RMR": common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}
}

func emptyRows() {
	mutex.Lock()
	defer mutex.Unlock()
	rows = map[string]common.Produce{}
}

// Allow sort by ProduceCode
type ByProduceCode []common.Produce

func (a ByProduceCode) Len() int      { return len(a) }
func (a ByProduceCode) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByProduceCode) Less(i, j int) bool {
	return strings.ToUpper(a[i].ProduceCode) < strings.ToUpper(a[j].ProduceCode)
}

// Helper function to test rows
func verifyRows(t *testing.T, count int, rows []common.Produce, expected []common.Produce) {
	sort.Sort(ByProduceCode(rows))
	sort.Sort(ByProduceCode(expected))

	if len(rows) != len(expected) || len(rows) != count {
		t.Errorf("ERROR -- lengths do not match:  expected(%+v) rows(%+v)  \n", expected, rows)
	}
	for i, r := range rows {
		if r.ProduceCode != expected[i].ProduceCode ||
			r.Name != expected[i].Name ||
			r.UnitPrice != expected[i].UnitPrice {
			t.Errorf("ERROR -- expected(%+v) found (%+v) \n", expected[i], r)
		}
	}
}

// Tests fetching all rows
func TestFetch(t *testing.T) {
	resetRows()
	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}

	var outputChannel chan common.Result
	var prows []common.Produce

	// Make sure all the rows look right
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	prows = []common.Produce{}
	for p := range outputChannel {
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		prows = append(prows, p.Prod)
	}

	verifyRows(t, len(expected), prows, expected)

	// What happens if the database is empty
	emptyRows()
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	for p := range outputChannel {
		if p.Err != "Row not found" || p.Count != 0 {
			t.Errorf("ERROR -- p(%v) does not have 'Row not found' for Err or Count is not 0\n", p)
		}
	}
}

// Tests fetching by Produce Code
func TestFetchByProduceCode(t *testing.T) {
	resetRows()
	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
	}

	outputChannel := make(chan common.Result, 2)
	go FetchByProduceCode("E5T6-9ui3-TH15-QR88", outputChannel)

	var rows []common.Produce
	for p := range outputChannel {
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		rows = append(rows, p.Prod)
	}

	verifyRows(t, len(expected), rows, expected)

	// Test an error case
	expected = []common.Produce{common.Produce{}}

	outputChannel = make(chan common.Result, 2)
	go FetchByProduceCode("A5T6-9ui3-TH15-QR88", outputChannel)

	rows = []common.Produce{}
	for p := range outputChannel {
		if p.Err != "Row not found" || p.Count != 0 {
			t.Errorf("ERROR -- p(%v) does not have 'Row not found' for Err or Count is not 0\n", p)
		}
		rows = append(rows, p.Prod)
	}

	verifyRows(t, len(expected), rows, expected)
}

// Tests adding rows
func TestAdd(t *testing.T) {
	resetRows()
	var addRows []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "ABCD-1234-ABCD-1234", Name: "Hamburger", UnitPrice: "505.460"},
		common.Produce{ProduceCode: "EFGH-2345-EFGH-2345", Name: "HotDogs", UnitPrice: "-3.46"},
		common.Produce{ProduceCode: "JKLM-5678-JKLM-5678", Name: "Buns", UnitPrice: "12.01"},
	}
	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
		common.Produce{ProduceCode: "ABCD-1234-ABCD-1234", Name: "Hamburger", UnitPrice: "505.460"},
		common.Produce{ProduceCode: "EFGH-2345-EFGH-2345", Name: "HotDogs", UnitPrice: "-3.46"},
		common.Produce{ProduceCode: "JKLM-5678-JKLM-5678", Name: "Buns", UnitPrice: "12.01"},
	}

	outputChannel := make(chan common.Result, 2)
	for _, p := range addRows {
		go Add(p, outputChannel)
	}

	var rrows []common.Result
	for _, _ = range addRows {
		rrows = append(rrows, <-outputChannel)
	}

	if len(rrows) != 3 {
		t.Errorf("ERROR - expected 3 rows returned.  Got (%v)\n", rrows)
	}

	prod := []common.Produce{}
	for _, r := range rrows {
		if r.Err != "" || r.Count != 1 {
			t.Errorf("ERROR - expected Err to be empty or Count == 1. Got (%v)\n", r)
		}
		prod = append(prod, r.Prod)
	}

	verifyRows(t, 3, prod, addRows)

	// Make sure all the rows look right
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	prows := []common.Produce{}
	for p := range outputChannel {
		// Could be lowercase - fix
		p.Prod.ProduceCode = strings.ToUpper(p.Prod.ProduceCode)
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		prows = append(prows, p.Prod)
	}

	verifyRows(t, len(expected), prows, expected)
}

// Tests adding Duplicate Row
func TestAddDuplicateRow(t *testing.T) {
	resetRows()
	var addRows []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "ABCD-1234-ABCD-1234", Name: "Hamburger", UnitPrice: "505.460"},
		common.Produce{ProduceCode: "abcd-1234-ABCD-1234", Name: "Hamburger", UnitPrice: "505.460"},
	}
	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
		common.Produce{ProduceCode: "ABCD-1234-ABCD-1234", Name: "Hamburger", UnitPrice: "505.460"},
	}
	expectedError := "ABCD-1234-ABCD-1234 ALREADY EXISTS" // Automatically capitialized on add

	outputChannel := make(chan common.Result, 2)
	for _, p := range addRows {
		go Add(p, outputChannel)
	}

	var rrows []common.Result
	for _, _ = range addRows {
		rrows = append(rrows, <-outputChannel)
	}

	if len(rrows) != 2 {
		t.Errorf("ERROR - expected 2 rows returned.  Got (%v)\n", rrows)
	}

	var prod []common.Produce
	if rrows[0].Count == 1 {
		// We could've added the lower case version - so uppercase it so verify will work
		rrows[0].Prod.ProduceCode = strings.ToUpper(rrows[0].Prod.ProduceCode)
		prod = append(prod, rrows[0].Prod)
		if strings.ToUpper(rrows[1].Err) != expectedError {
			t.Errorf("ERROR - expected (%v) got received (%v)\n", expectedError, strings.ToUpper(rrows[1].Err))
		}
	} else {
		// We could've added the lower case version - so uppercase it so verify will work
		rrows[1].Prod.ProduceCode = strings.ToUpper(rrows[1].Prod.ProduceCode)
		prod = append(prod, rrows[1].Prod)
		if strings.ToUpper(rrows[0].Err) != expectedError {
			t.Errorf("ERROR - expected (%v) got received (%v)\n", expectedError, strings.ToUpper(rrows[0].Err))
		}
	}

	verifyRows(t, 1, prod, addRows[:1])

	// Make sure all the rows look right
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	prows := []common.Produce{}
	for p := range outputChannel {
		// Could be lowercase - fix
		p.Prod.ProduceCode = strings.ToUpper(p.Prod.ProduceCode)
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		prows = append(prows, p.Prod)
	}

	verifyRows(t, len(expected), prows, expected)
}

// Tests DeleteRow
func TestDeleteRow(t *testing.T) {
	resetRows()
	var delRow []common.Produce = []common.Produce{common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"}}

	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}

	outputChannel := make(chan common.Result, 2)
	for _, p := range delRow {
		go DeleteRow(p, outputChannel)
	}

	var rrows []common.Result
	for _, _ = range delRow {
		rrows = append(rrows, <-outputChannel)
	}

	if len(rrows) != 1 || rrows[0].Count != 1 || rrows[0].Err != "" {
		t.Errorf("ERROR - expected 1 rows returned with Count == 1 and Err string empty. Got (%v)\n", rrows)
	}

	prows := []common.Produce{rrows[0].Prod}

	verifyRows(t, 1, prows, delRow)

	// Make sure all the rows look right
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	prows = []common.Produce{}
	for p := range outputChannel {
		// Could be lowercase - fix
		p.Prod.ProduceCode = strings.ToUpper(p.Prod.ProduceCode)
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		prows = append(prows, p.Prod)
	}

	verifyRows(t, len(expected), prows, expected)
}

// Tests DeleteRow with a bad ProduceCode (ProduceRow not found)
func TestDeleteRowBadProduceCode(t *testing.T) {
	resetRows()
	var delRow []common.Produce = []common.Produce{common.Produce{ProduceCode: "ABC6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"}}

	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}

	var expected_error = "Row not found"

	outputChannel := make(chan common.Result, 2)
	for _, p := range delRow {
		go DeleteRow(p, outputChannel)
	}

	var rrows []common.Result
	for _, _ = range delRow {
		rrows = append(rrows, <-outputChannel)
	}

	if len(rrows) != 1 || rrows[0].Count != 0 || rrows[0].Err != expected_error {
		t.Errorf("ERROR - expected 1 rows returned with Count == 0 and Err string (%v). Got (%v)\n", expected_error, rrows)
	}

	prows := []common.Produce{rrows[0].Prod}

	verifyRows(t, 1, prows, []common.Produce{common.Produce{}})

	// Make sure all the rows look right
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	prows = []common.Produce{}
	for p := range outputChannel {
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		prows = append(prows, p.Prod)
	}

	verifyRows(t, len(expected), prows, expected)
}

// Tests Deleting a row
func TestDelete(t *testing.T) {
	resetRows()
	var delRow []common.Produce = []common.Produce{common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"}}

	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}

	outputChannel := make(chan common.Result, 2)
	for _, p := range delRow {
		go Delete(p.ProduceCode, outputChannel)
	}

	var rrows []common.Result
	for _, _ = range delRow {
		rrows = append(rrows, <-outputChannel)
	}

	if len(rrows) != 1 || rrows[0].Count != 1 || rrows[0].Err != "" {
		t.Errorf("ERROR - expected 1 rows returned with Count == 1 and Err string empty. Got (%v)\n", rrows)
	}

	prows := []common.Produce{rrows[0].Prod}

	verifyRows(t, 1, prows, delRow)

	// Make sure all the rows look right
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	prows = []common.Produce{}
	for p := range outputChannel {
		// Could be lowercase - fix
		p.Prod.ProduceCode = strings.ToUpper(p.Prod.ProduceCode)
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		prows = append(prows, p.Prod)
	}

	verifyRows(t, len(expected), prows, expected)
}

// Tests Deleting a bad ProduceCode (Produce not found)
func TestDeleteBadProduceCode(t *testing.T) {
	resetRows()
	var delRow []common.Produce = []common.Produce{common.Produce{ProduceCode: "ABC6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"}}

	var expected []common.Produce = []common.Produce{
		common.Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.46"},
		common.Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		common.Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		common.Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}

	var expected_error = "Row not found"

	outputChannel := make(chan common.Result, 2)
	for _, p := range delRow {
		go Delete(p.ProduceCode, outputChannel)
	}

	var rrows []common.Result
	for _, _ = range delRow {
		rrows = append(rrows, <-outputChannel)
	}

	if len(rrows) != 1 || rrows[0].Count != 0 || rrows[0].Err != expected_error {
		t.Errorf("ERROR - expected 1 rows returned with Count == 0 and Err string (%v). Got (%v)\n", expected_error, rrows)
	}

	prows := []common.Produce{rrows[0].Prod}

	verifyRows(t, 1, prows, []common.Produce{common.Produce{}})

	// Make sure all the rows look right
	outputChannel = make(chan common.Result, 2)
	go Fetch(outputChannel)

	prows = []common.Produce{}
	for p := range outputChannel {
		if p.Err != "" || p.Count != 1 {
			t.Errorf("ERROR -- p(%v) does not have nil for Err or Count is not 1\n", p)
		}
		prows = append(prows, p.Prod)
	}

	verifyRows(t, len(expected), prows, expected)
}
