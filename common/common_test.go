package common

import (
	"testing"
)

// Helper function to test
func verify(t *testing.T, function string, expected bool, input string, result bool) {
	if expected != result {
		t.Errorf("ERROR - for (%v) expected (%v) for input (%v) but got (%v)\n", function, expected, input, result)
	}
}

// validateProduceCodeTestStruct
type vpcTS struct {
	exp   bool   // Expected
	input string // Input
}

// validateProduceCodeTestStructs: test cases
var vpcTSs = []vpcTS{
	{false, ""},
	{true, "A12T-4Gh7-QPL9-3N4M"},
	{false, "A12T-4Gh7-QPL9-3N4MB"},
	{false, "A12T_4Gh7-QPL9-3N4M"},
	{false, "A12T!4Gh7-QPL9-3N4M"},
	{false, " A12T-4Gh7-QPL9-3N4M"},
	{false, "A12T-4Gh7-QPL9-3N4M "},
	{false, " A12T-4Gh7-QPL9-3N4M "},
}

// Verify validate Produce Code RegEx
func TestValidateProduceCode(t *testing.T) {
	for _, tt := range vpcTSs {
		verify(t, "validateProduceCode", tt.exp, tt.input, ValidateProduceCode(tt.input))
	}
}

// validateNameRegExTestStruct
type vnrxTS struct {
	exp   bool   // Expected
	input string // Input
}

// validateNameRegExTestStructs: test cases
var vnrxTSs = []vnrxTS{
	{false, ""},
	{true, "Pe4ch"},
	{true, "Pe4ch Pits"},
	{false, "Pe4ch "},
	{false, " Pe4ch"},
	{false, "Pe!ch"},
}

// Verify validate Name RegEx
func TestValidateNameRegEx(t *testing.T) {
	for _, tt := range vnrxTSs {
		verify(t, "validateName", tt.exp, tt.input, validateName(tt.input))
	}
}

// validateUnitPriceTestStruct
type vupTS struct {
	exp   bool   // Expected
	input string // Input
}

// validateUnitPriceTestStructs;: test cases
var vupTSs = []vupTS{
	{false, ""},
	{false, ""},
	{true, "."},
	{true, "3.46"},
	{true, "$3.46"},
	{true, "$3.4"},
	{false, "$A.46"},
	{false, "A.46"},
	{false, "3.A6"},
	{false, "3.460"},
	{false, "3.46 "},
	{false, " 3.46"},
}

// Verify validate UnitPrice RegEx
func TestValidateUnitPrice(t *testing.T) {
	for _, tt := range vupTSs {
		verify(t, "validateUnitPrice", tt.exp, tt.input, validateUnitPrice(tt.input))
	}
}

// Helper function to test
func verifyfixUnitPrice(t *testing.T, function string, expected string, input string, result string) {
	if expected != result {
		t.Errorf("ERROR - for (%v) expected (%v) for input (%v) but got (%v)\n", function, expected, input, result)
	}
}

// validateFixUnitPriceTestStruct
type vfupTS struct {
	function string
	expected string
	input    string
}

// validateFixUnitPriceTestStructs : test cases
var vfupTSs = []vfupTS{
	{"fixUnitPrice", "0.55", "0.55"},
	{"fixUnitPrice", "0.55", "$0.55"},
	{"fixUnitPrice", "0.00", "$0."},
	{"fixUnitPrice", "0.00", "$."},
	{"fixUnitPrice", "0.00", "."},
	{"fixUnitPrice", "0.00", ".0"},
	{"fixUnitPrice", "0.00", ".00"},
}

// Verify fixUnitPrice
func TestFixUnitPrice(t *testing.T) {
	for _, tt := range vfupTSs {
		verifyfixUnitPrice(t, tt.function, tt.expected, tt.input, fixUnitPrice(tt.input))
	}
}

// Helper function to validate ValidateProduce
func verifyValidateProduce(t *testing.T, function string, expected bool, input Produce, result bool) {
	if expected != result {
		t.Errorf("ERROR - for (%v) expected (%v) for input (%v) but got (%v)\n", function, expected, input, result)
	}
}

// helper function to validate the errors from ValidateProduce
func verifyValidateProduceErrors(t *testing.T, function string, expectedErrors []string, resultErrors []string) {
	if len(expectedErrors) != len(resultErrors) {
		t.Errorf("ERROR - for (%v) expectedErrors length (%v) did not match resultErrors length (%v)\n", function, len(expectedErrors), len(resultErrors))
	}

	for i, _ := range expectedErrors {
		if expectedErrors[i] != resultErrors[i] {
			t.Errorf("ERROR - for (%v) expectedErrors[%v](%v) did not match resultErrors[%v](%v)\n", function, i, expectedErrors[i], i, resultErrors[i])
		}
	}
}

// validateProduceTestStruct
type vpTS struct {
	function       string
	input          Produce
	expected       bool
	expectedErrors []string
}

// validateProduceTestStructs : test cases
var vpTSs = []vpTS{
	{function: "ValidateProduce",
		input:          Produce{ProduceCode: "", Name: "", UnitPrice: ""},
		expected:       false,
		expectedErrors: []string{"Detected error for Produce Code ()", "Detected error for Produce Name ()", "Detected error for Produce Unit Price ()"}},
	{function: "ValidateProduce",
		input:          Produce{ProduceCode: "ABCD-ABCD-ABCD-ABCD", Name: "TEST", UnitPrice: "."},
		expected:       true,
		expectedErrors: []string{}},
}

// Verify ValidateProduce
func TestValidateProduce(t *testing.T) {
	var result bool
	var resultErrors []string

	for _, tt := range vpTSs {
		result, resultErrors = ValidateProduce(tt.input)
		verifyValidateProduce(t, tt.function, tt.expected, tt.input, result)
		verifyValidateProduceErrors(t, tt.function, tt.expectedErrors, resultErrors)
	}
}

// helper function to validate the errors from FixProduce
func verifyFixProduce(t *testing.T, function string, expected Produce, input Produce, result Produce) {
	if expected != result {
		t.Errorf("ERROR - for (%v) expected unitPrice = (%v)  input unitPrice = (%v) result unitPrice = (%v) \n", function, expected.UnitPrice, input.UnitPrice, result.UnitPrice)
	}
}

// fixProduceTestStruct
type fpTS struct {
	function string
	expected Produce
	input    Produce
}

// fixProduceTestStructs: test cases
var fpTSs = []fpTS{
	{function: "fixProduce",
		expected: Produce{ProduceCode: "ABCD-ABCD-ABCD-ABCD", Name: "TEST", UnitPrice: "0.00"},
		input:    Produce{ProduceCode: "ABCD-ABCD-ABCD-ABCD", Name: "TEST", UnitPrice: "."}},
}

func TestFixProduce(t *testing.T) {
	for _, tt := range fpTSs {
		verifyFixProduce(t, tt.function, tt.expected, tt.input, FixProduce(tt.input))
	}
}
