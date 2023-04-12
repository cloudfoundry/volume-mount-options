package utils_test

import (
	"code.cloudfoundry.org/volume-mount-options/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("KernelMountOptions", func() {
	Describe("#ParseOptionStringToMap", func() {
		var (
			optionString string
			opts         map[string]interface{}
		)

		BeforeEach(func() {
			optionString = ""
		})

		JustBeforeEach(func() {
			opts = utils.ParseOptionStringToMap(optionString, "=")
		})

		Context("given an empty option string", func() {
			It("should return an empty map of options", func() {
				Expect(opts).To(BeEmpty())
			})
		})

		Context("given an option string", func() {
			BeforeEach(func() {
				optionString = "opt1=val1,opt2=val2"
			})

			It("should return an map of options", func() {
				Expect(opts).To(Equal(map[string]interface{}{
					"opt1": "val1",
					"opt2": "val2",
				}))
			})
		})

		Context("given an option without a value", func() {
			BeforeEach(func() {
				optionString = "opt1=val1,opt2"
			})

			It("should return an map of options", func() {
				Expect(opts).To(Equal(map[string]interface{}{
					"opt1": "val1",
					"opt2": "",
				}))
			})
		})

		Context("given an option value that includes a equal sign", func() {
			BeforeEach(func() {
				optionString = "opt1=val1,opt2=val2=val3"
			})

			It("should return an map of options", func() {
				Expect(opts).To(Equal(map[string]interface{}{
					"opt1": "val1",
					"opt2": "val2=val3",
				}))
			})
		})
	})
})
