package utils_test

import (
	"fmt"
	"math"
	"strconv"

	"code.cloudfoundry.org/volume-mount-options/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeConversion", func() {
	var (
		stringInput  string
		int64Input   int64
		float64Input float64
		boolInput    bool
		invalidInput func()
	)

	stringInput = "a string"
	int64Input = math.MaxInt64
	float64Input = math.MaxFloat64
	boolInput = true
	invalidInput = func() {}

	DescribeTable(
		"#InterfaceToString",
		func(input interface{}, expected string) {
			output := utils.InterfaceToString(input)
			Expect(output).To(Equal(expected))
		},
		Entry("string input", stringInput, "a string"),
		Entry("int64 input", int64Input, fmt.Sprintf("%d", int64Input)),
		Entry("float64 input", float64Input, fmt.Sprintf("%f", float64Input)),
		Entry("bool input", boolInput, strconv.FormatBool(boolInput)),
		Entry("invalid input", invalidInput, ""),
	)
})
