package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func getEcho() *echo.Echo {
	e := echo.New()

	// Add a new Produce item to Inventory
	e.POST("/produce", AddProduce)

	// Delete Produce item from Inventory
	e.DELETE("/produce/:ProduceCode", DeleteProduce)

	// Fetch all Produce items from Inventory
	e.GET("/produce", FetchProduce)

	// Fetch a Produce item from Inventory by Produce Code
	e.GET("/produce/:ProduceCode", FetchProduceByProduceCode)

	return e
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

// Test FetchProduce
func TestFetchProduce(t *testing.T) {
	// Success condition
	expected := http.StatusOK

	req := httptest.NewRequest(echo.GET, "/produce", nil)
	rec := httptest.NewRecorder()

	e := getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestFetchProduce** - Success - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)
}

// Test FetchProduceByProduceCode
func TestFetchProduceByProduceCode(t *testing.T) {
	// Success condition
	expected := http.StatusOK
	req := httptest.NewRequest(echo.GET, "/produce/A12T-4gh7-QPL9-3N4M", nil)
	rec := httptest.NewRecorder()

	e := getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestFetchProduceByProduceCode** - Success - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)

	// Produce Not found
	expected = http.StatusNoContent
	req = httptest.NewRequest(echo.GET, "/produce/ZZZZ-4gh7-QPL9-3N4M", nil)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestFetchProduceByProduceCode** - Produce not found - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)

	// Produce Code misformed
	expected = http.StatusBadRequest
	req = httptest.NewRequest(echo.GET, "/produce/-ZZZZ-4gh7-QPL9-3N4M", nil)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestFetchProduceByProduceCode** - Produce Code misformed - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)
}

// Test AddProduct
func TestAddProduce(t *testing.T) {
	// Success condition
	expected := http.StatusOK
	expectedBody := "{\"Produce\":[{\"Produce Code\":\"AAAA-1111-2222-3333\",\"Name\":\"Pizza Pie\",\"Unit Price\":\"200.6\"}]}"
	req := httptest.NewRequest(echo.POST, "/produce", strings.NewReader(" [ {\"Produce Code\": \"AAAA-1111-2222-3333\", \"Name\": \"Pizza Pie\", \"Unit Price\": \"200.6\" } ]"))
	rec := httptest.NewRecorder()

	e := getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected(%v) received(%v) expectedBody(%v) receivedBody (%v) \n", expected, rec.Code, expectedBody, rec.Body)
	}

	// Failure condition - bad body
	expected = http.StatusBadRequest
	expectedBody = "{\"Rejected Produce\":[{\"Errors\":[\"Failed to read request body\"]}]}"
	req = httptest.NewRequest(echo.POST, "/produce", errReader(0))
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected(%v) received(%v) expectedBody(%v) receivedBody (%v) \n", expected, rec.Code, expectedBody, rec.Body)
	}

	// Failure condition - bad JSON
	expected = http.StatusBadRequest
	expectedBody = "{\"Rejected Produce\":[{\"Errors\":[\"Failed to unmarshal request body\"]}]}"
	req = httptest.NewRequest(echo.POST, "/produce", strings.NewReader(" {\"Produce Code\": \"BBBB-1111-2222-3333\", \"Name\": \"Pizza Pie\", \"Unit Price\": \"200.6\"  "))
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected(%v) received(%v) expectedBody(%v) receivedBody (%v) \n", expected, rec.Code, expectedBody, rec.Body)
	}

	// Failure condition - invalid Produce form
	expected = http.StatusBadRequest
	expectedBody = "{\"Rejected Produce\":[{\"Produce\":{\"Produce Code\":\"BBBB-1111-2222-3333-\",\"Name\":\" Pizza Pie \",\"Unit Price\":\"200.645\"},\"Errors\":[\"Detected error for Produce Code (BBBB-1111-2222-3333-)\",\"Detected error for Produce Name ( Pizza Pie )\",\"Detected error for Produce Unit Price (200.645)\"]}]}"
	req = httptest.NewRequest(echo.POST, "/produce", strings.NewReader(" {\"Produce Code\": \"BBBB-1111-2222-3333-\", \"Name\": \" Pizza Pie \", \"Unit Price\": \"200.645\" } "))
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected(%v) received(%v) expectedBody(%v) receivedBody (%v) \n", expected, rec.Code, expectedBody, rec.Body)
	}

	// Failure condition - duplicate Produce Code
	expected = http.StatusPartialContent
	expectedBody = "{\"Rejected Produce\":[{\"Produce\":{\"Produce Code\":\"AAAA-1111-2222-3333\",\"Name\":\"Pizza Pie\",\"Unit Price\":\"200.6\"},\"Errors\":[\"AAAA-1111-2222-3333 already exists\"]}]}"
	req = httptest.NewRequest(echo.POST, "/produce", strings.NewReader(" [ {\"Produce Code\": \"AAAA-1111-2222-3333\", \"Name\": \"Pizza Pie\", \"Unit Price\": \"200.6\" } ]"))
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected(%v) received(%v) expectedBody(%v) receivedBody (%v) \n", expected, rec.Code, expectedBody, rec.Body)
	}

	// Failure condition - duplicate Produce Code with add
	expected = http.StatusPartialContent
	expectedBody = "{\"Produce\":[{\"Produce Code\":\"AAAA-1111-2222-9999\",\"Name\":\"Black Truffles\",\"Unit Price\":\"200.6\"}],\"Rejected Produce\":[{\"Produce\":{\"Produce Code\":\"AAAA-1111-2222-3333\",\"Name\":\"Pizza Pie\",\"Unit Price\":\"200.6\"},\"Errors\":[\"AAAA-1111-2222-3333 already exists\"]}]}"
	req = httptest.NewRequest(echo.POST, "/produce", strings.NewReader(" [ {\"Produce Code\": \"AAAA-1111-2222-3333\", \"Name\": \"Pizza Pie\", \"Unit Price\": \"200.6\" }, {\"Produce Code\": \"AAAA-1111-2222-9999\", \"Name\": \"Black Truffles\", \"Unit Price\": \"200.6\" } ]"))
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected(%v) received(%v) expectedBody(%v) receivedBody (%v) \n", expected, rec.Code, expectedBody, rec.Body)
	}

	// Failure condition - Produce Code with add and failure
	expected = http.StatusPartialContent
	expectedBody = "{\"Produce\":[{\"Produce Code\":\"AAAA-1111-2222-7777\",\"Name\":\"Black Truffles\",\"Unit Price\":\"200.6\"}],\"Rejected Produce\":[{\"Produce\":{\"Produce Code\":\"-AAAA-1111-2222-3333\",\"Name\":\"Pizza Pie\",\"Unit Price\":\"200.6\"},\"Errors\":[\"Detected error for Produce Code (-AAAA-1111-2222-3333)\"]}]}"
	req = httptest.NewRequest(echo.POST, "/produce", strings.NewReader(" [ {\"Produce Code\": \"-AAAA-1111-2222-3333\", \"Name\": \"Pizza Pie\", \"Unit Price\": \"200.6\" }, {\"Produce Code\": \"AAAA-1111-2222-7777\", \"Name\": \"Black Truffles\", \"Unit Price\": \"200.6\" } ]"))
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected(%v) received(%v) expectedBody(%v) receivedBody (%v) \n", expected, rec.Code, expectedBody, rec.Body)
	}
}

// Test Delete Produce
func TestDeleteProduce(t *testing.T) {
	// Success condition
	expected := http.StatusOK
	req := httptest.NewRequest(echo.DELETE, "/produce/A12T-4gh7-QPL9-3N4M", nil)
	rec := httptest.NewRecorder()

	e := getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestDeleteProduce** - Success - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)

	// Failure condition, Produce not found
	expected = http.StatusNotFound
	req = httptest.NewRequest(echo.DELETE, "/produce/A12T-4gh7-QPL9-3N4M", nil)
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestDeleteProduce** - Produce not found - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)

	// Failure condition, Produce Code misformed
	expected = http.StatusBadRequest
	req = httptest.NewRequest(echo.DELETE, "/produce/-A12T-4gh7-QPL9-3N4M", nil)
	rec = httptest.NewRecorder()

	e = getEcho()
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestDeleteProduce** - failure Condition - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)
}

// Test FetchProduceOnEmptyDB
func TestFetchProduceOnEmptyDB(t *testing.T) {
	expected := http.StatusOK

	req := httptest.NewRequest(echo.GET, "/produce", nil)
	rec := httptest.NewRecorder()

	e := getEcho()
	e.ServeHTTP(rec, req)

	// Failure - no row found
	// Must delete all the rows
	req = httptest.NewRequest(echo.DELETE, "/produce/A12T-4GH7-QPL9-3N4M", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/E5T6-9UI3-TH15-QR88", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/YRT6-72AS-K736-L4AR", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/TQ4C-VV6T-75ZX-1RMR", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/AAAA-1111-2222-3333", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/BBBB-1111-2222-3333", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/BBBB-1111-2222-3333", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/AAAA-1111-2222-7777", nil)
	e.ServeHTTP(rec, req)
	req = httptest.NewRequest(echo.DELETE, "/produce/AAAA-1111-2222-9999", nil)
	e.ServeHTTP(rec, req)

	expected = http.StatusNoContent

	req = httptest.NewRequest(echo.GET, "/produce", nil)
	rec = httptest.NewRecorder()

	//e := getEcho() // Don't want to reset echo
	e.ServeHTTP(rec, req)

	if expected != rec.Code {
		t.Errorf("ERROR -- expected (%v) but got (%v) body is (%v) \n", expected, rec.Code, rec.Body)
	}
	log.Printf("**TestFetchProduceOnEmptyDB** - Status is (%v) Body is (%v)\n", rec.Code, rec.Body)
}
