package utils_test

import (
	"fmt"
	"math"

	"code.cloudfoundry.org/volume-mount-options/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KernelMountOptions", func() {
	Describe("#ToKernelMountOptionString", func() {
		var (
			mountOpts          map[string]interface{}
			kernelMountOptions string
		)

		BeforeEach(func() {
			mountOpts = make(map[string]interface{})
		})

		JustBeforeEach(func() {
			kernelMountOptions = utils.ToKernelMountOptionString(mountOpts)
		})

		Context("given an empty mount opts", func() {
			It("should return an empty mount opts string", func() {
				Expect(kernelMountOptions).To(BeEmpty())
			})
		})

		Context("given a mount opts", func() {
			BeforeEach(func() {
				mountOpts = map[string]interface{}{
					"opt1": "val1",
					"opt2": "val2",
				}
			})

			It("should return a valid mount opts string", func() {
				Expect(kernelMountOptions).To(Equal("opt1=val1,opt2=val2"))
			})
		})

		Context("given an integer option value with a leading zero", func() {
			BeforeEach(func() {
				mountOpts = map[string]interface{}{
					"opt1": "0123",
				}
			})

			It("strips the leading zero from the mount option string", func() {
				Expect(kernelMountOptions).To(Equal("opt1=123"))
			})
		})

		Context("given an integer option value", func() {
			BeforeEach(func() {
				mountOpts = map[string]interface{}{
					"opt1": math.MaxInt64,
				}
			})

			It("strips the leading zero from the mount option string", func() {
				Expect(kernelMountOptions).To(Equal(fmt.Sprintf("opt1=%d", math.MaxInt64)))
			})
		})

		Context("given a mount option with no value", func() {
			BeforeEach(func() {
				mountOpts = map[string]interface{}{
					"does-not-matter": "",
				}
			})

			It("adds the mount option to the string without a value", func() {
				Expect(kernelMountOptions).To(Equal("does-not-matter"))
			})
		})
	})

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
