package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"example.com/produce_demo/common"
	"example.com/produce_demo/db"

	"github.com/labstack/echo/v4"
)

// FetchMsg return structure - used by FetchProduce and FetchProduceByProduceCode
type FetchMsg struct {
	Err     string            `json:"Error,omitempty"`
	Produce *[]common.Produce `json:"Produce,omitempty"`
}

// Fetch all Produce concurrently
func FetchProduce(c echo.Context) error {

	// Fetch rows
	outputChannel := make(chan common.Result, 1)
	go db.Fetch(outputChannel)

	// Process Results
	errorString := ""
	produceList := []common.Produce{}
	for p := range outputChannel {
		if p.Count != 1 || p.Err != "" {
			errorString = errorString + p.Err
		} else {
			produceList = append(produceList, p.Prod)
		}
	}

	log.Printf("FetchProduce - produceList is (%v)\n", produceList)

	// Handle No rows found
	if len(produceList) == 0 {
		return c.JSON(http.StatusNoContent, FetchMsg{Err: "No produce found"}) // Returns 204
	}
	// Handle Errors
	if errorString != "" { // Unreachable code - fetch only sets errorString if no rows returned
		return c.JSON(http.StatusInternalServerError, FetchMsg{Err: "Internal Error detected"}) // Returns 500
	}

	// Final Return
	return c.JSON(http.StatusOK, FetchMsg{Produce: &produceList}) // Returns 200
}

// Fetch Produce by ProduceCode
func FetchProduceByProduceCode(c echo.Context) error {

	// Get and Validate Param
	produceCode := c.Param("ProduceCode")
	if !common.ValidateProduceCode(produceCode) {
		log.Printf("FetchProduceByProduceCode - failed with produceCode(%v)\n", produceCode)
		return c.JSON(http.StatusBadRequest, FetchMsg{Err: "Bad Produce Code"}) // Returns 400
	}

	// Fetch Rows
	outputChannel := make(chan common.Result, 1)
	go db.FetchByProduceCode(produceCode, outputChannel)

	// Process Results
	errorString := ""
	produceList := []common.Produce{}
	for p := range outputChannel {
		if p.Count != 1 || p.Err != "" {
			errorString = errorString + p.Err
		} else {
			produceList = append(produceList, p.Prod)
		}
	}

	log.Printf("FetchProduceByProduceCode - produceList is (%v)\n", produceList)

	// Handle Errors
	if len(produceList) == 0 {
		return c.JSON(http.StatusNoContent, FetchMsg{Err: "No produce found"}) // Returns 204
	}
	if errorString != "" { // Unreachable code - fetch only sets errorString if no rows returned
		return c.JSON(http.StatusInternalServerError, FetchMsg{Err: "Internal Error detected"}) // Returns 500
	}

	// Final Return
	return c.JSON(http.StatusOK, FetchMsg{Produce: &produceList}) // Returns 200
}

// Return Structures for addProduceCall - used for both success and failure conditions
type ErrorProduce struct {
	Produce *common.Produce `json:"Produce,omitempty"`
	Errors  []string        `json:"Errors"`
}

type ReturnAdd struct {
	Produce         []common.Produce `json:"Produce,omitempty"`
	RejectedProduce []ErrorProduce   `json:"Rejected Produce,omitempty"`
}

// For a list of produce - run common.ValidateProduce on each:
//   Valid produce is added to the validProduceList
//   Invalid produce is added (along with error) to RejectedProduce
func getValidProduceList(produceList []common.Produce) ([]common.Produce, []ErrorProduce) {
	rejectedProduce := []ErrorProduce{}
	validProduceList := []common.Produce{}

	// Verify productList
	for i, p := range produceList {
		ret, validProduceError := common.ValidateProduce(p)
		if ret != true {
			rejectedProduce = append(rejectedProduce, ErrorProduce{Produce: &produceList[i], Errors: validProduceError}) // NOTE: subtle error if using &p
		} else {
			validProduceList = append(validProduceList, p)
		}
	}
	return validProduceList, rejectedProduce
}

// Add Produce Concurrently
func AddProduce(c echo.Context) error {
	defer c.Request().Body.Close()

	var produceList []common.Produce // NOTE: Here we just need a variable to bind to

	// Ready the body of the POST - fail if we can't read it
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("AddProduce - Failed reading the request body for AddProduce: %s\n", err)
		errProduce := ErrorProduce{Errors: []string{"Failed to read request body"}}
		return c.JSON(http.StatusBadRequest, ReturnAdd{RejectedProduce: []ErrorProduce{errProduce}}) // Returns 400
	}

	log.Printf("AddProduce - Body is (%v)\n", string(b))

	// Unmarshal the body into a productList
	err = json.Unmarshal(b, &produceList)
	if err != nil {
		// if unmarshal of list fails - try unmarshal of just Produce
		var produce common.Produce // NOTE: Here we just need a variable to bind to

		err = json.Unmarshal(b, &produce)
		if err != nil {
			// Okay, give up.
			log.Printf("AddProduce - Failed unmarshalling in AddProduce: %s\n", err)
			errProduce := ErrorProduce{Errors: []string{"Failed to unmarshal request body"}}
			return c.JSON(http.StatusBadRequest, ReturnAdd{RejectedProduce: []ErrorProduce{errProduce}}) // Returns 400
		}
		produceList = append(produceList, produce)
	}

	// getValidProduce
	validProduceList, rejectedProduceList := getValidProduceList(produceList)

	// Attempt to add validProduce
	addedProduceList := []common.Produce{}
	if len(validProduceList) > 0 {
		outputChannel := make(chan common.Result, 2)
		for _, p := range validProduceList {
			go db.Add(p, outputChannel)
		}

		// Get the results
		var r common.Result
		for _, _ = range validProduceList {
			r = <-outputChannel
			if r.Err != "" {
				rejectedProduceList = append(rejectedProduceList, ErrorProduce{Produce: &r.Prod, Errors: []string{r.Err}})
			} else {
				addedProduceList = append(addedProduceList, r.Prod)
			}
		}

		// Handle Errors
		if len(rejectedProduceList) != 0 {
			if len(addedProduceList) != 0 {
				return c.JSON(http.StatusPartialContent, ReturnAdd{Produce: addedProduceList, RejectedProduce: rejectedProduceList}) //Returns 206
			} else {
				return c.JSON(http.StatusPartialContent, ReturnAdd{RejectedProduce: rejectedProduceList}) //Returns 206
			}
		}
	}

	// Handle Errors - but we never called db.Add
	if len(rejectedProduceList) != 0 {
		return c.JSON(http.StatusBadRequest, ReturnAdd{RejectedProduce: rejectedProduceList}) // Returns 400
	}

	// Final Return
	return c.JSON(http.StatusOK, ReturnAdd{Produce: addedProduceList}) // Returns 200
}

// DeleteReturn structure - used by DeleteProduce
type DeleteReturn struct {
	Msg string `json:"Msg,omitempty"`
	Err string `json:"Error,omitempty"`
}

// Delete Produce by ProduceCode concurrently
func DeleteProduce(c echo.Context) error {

	// Get and Validate Param
	produceCode := c.Param("ProduceCode")
	if !common.ValidateProduceCode(produceCode) {
		log.Printf("DeleteProduce - failed with produceCode(%v)\n", produceCode)
		return c.JSON(http.StatusBadRequest, DeleteReturn{Err: "Bad Produce Code"}) // Return 400
	}

	// Delete Row
	outputChannel := make(chan common.Result, 2)
	go db.Delete(produceCode, outputChannel)

	// Get the results
	var r common.Result
	r = <-outputChannel

	// Handle Errors
	if r.Err != "" {
		log.Printf("DeleteProduce - Detected Error (%s)\n", r.Err)
		return c.JSON(http.StatusNotFound, DeleteReturn{Err: "Produce not found"}) // Returns 404
	}

	// Final Return
	return c.JSON(http.StatusOK, DeleteReturn{Msg: "Produce " + produceCode + " deleted"}) // Returns 200
}
