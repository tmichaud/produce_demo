// Common Structures, Validation, and Correction (Fix) logic
package common

import (
	"log"
	"regexp"
)

// Common struct and variables
// Since persistence layer and model match, we are using a single structure
// If we want seperation of model and persistence - conversion method could be placed here

// type Row struct {
//	ProduceCode string `json:"Produce Code"`
//	Name        string `json:"Name"`
//	UnitPrice   string `json:"Unit Price"`
//}

// NOTE: Considered using github.com/shopspring/decimal for UnitPrice.
// NOTE: Seems like overkill for this at the moment. Will refactor if necessary.
// NOTE: Obviously float/double are out since we are talking about monetary amount.
//       Using String at the moment, since we don't seem to be doing any arithmetic.

// Produce structure used for both api and db
type Produce struct {
	ProduceCode string `json:"Produce Code"`
	Name        string `json:"Name"`
	UnitPrice   string `json:"Unit Price"`
}

// Communication between api/handler and db
type Result struct {
	Prod  Produce
	Err   string
	Count int
}

// RegEx for Produce's Produce Codes
// The produce codes are sixteen characters long, with dashes separating each four character group
// The produce codes are alphanumeric and case insensitive
var validateProduceCodeRegEx string = "^[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}$"

// RegEx for Produce's Name
// The produce name is alphanumeric and case insensitive (and may contain spaces)
var validateNameRegEx string = "^[[:alnum:]]([[:alnum:]]| )+[[:alnum:]]$"

// RegEx for Produce's Unit Price
// The produce unit price is a number with up to 2 decimal places
// And may optionally start with a '$'
var validateUnitPriceRegEx string = "^[$]{0,1}[[:digit:]]*[.][[:digit:]]{0,2}$"

// tests if a Produce's Produce Code is valid
func ValidateProduceCode(produceCode string) bool {
	re := regexp.MustCompile(validateProduceCodeRegEx)
	return re.MatchString(produceCode)
}

// tests if a Produce's Name is valid
func validateName(name string) bool {
	re := regexp.MustCompile(validateNameRegEx)
	return re.MatchString(name)
}

// tests if Produce's Unit Price is valid
func validateUnitPrice(unitPrice string) bool {
	re := regexp.MustCompile(validateUnitPriceRegEx)
	return re.MatchString(unitPrice)
}

// fixes a Produce's Unit Price by removing leading '$' and adding '0' before and up to '00' after decimal
func fixUnitPrice(unitPrice string) string {
	if '$' == unitPrice[0] {
		unitPrice = unitPrice[1:]
	}
	if '.' == unitPrice[0] {
		unitPrice = "0" + unitPrice
	}
	if '.' == unitPrice[len(unitPrice)-1] {
		unitPrice = unitPrice + "00"
	}
	if '.' == unitPrice[len(unitPrice)-2] {
		unitPrice = unitPrice + "0"
	}
	return unitPrice
}

// Convenience Method to test all fields of a Produce
func ValidateProduce(p Produce) (bool, []string) {
	ret := true
	errorText := []string{}
	if ValidateProduceCode(p.ProduceCode) != true {
		log.Printf("ValidateProduceCode failed for produce(%v)", p)
		errorText = append(errorText, "Detected error for Produce Code ("+p.ProduceCode+")")
		ret = false
	}
	if validateName(p.Name) != true {
		log.Printf("ValidateName failed for produce(%v)", p)
		errorText = append(errorText, "Detected error for Produce Name ("+p.Name+")")
		ret = false
	}
	if validateUnitPrice(p.UnitPrice) != true {
		log.Printf("ValidateUnitPrice failed for produce(%v)", p)
		errorText = append(errorText, "Detected error for Produce Unit Price ("+p.UnitPrice+")")
		ret = false
	}
	return ret, errorText
}

// Convenience Method to fix all fields of a Produce
func FixProduce(p Produce) Produce {
	p.UnitPrice = fixUnitPrice(p.UnitPrice)
	return p
}
