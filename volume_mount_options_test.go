package volume_mount_options_test

import (
	"errors"
	"fmt"

	vmo "code.cloudfoundry.org/volume-mount-options"
	volumemountoptionsfakes "code.cloudfoundry.org/volume-mount-options/volume-mount-optionsfakes"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeMountOptions", func() {
	Describe("#NewMountOpts", func() {
		var (
			allowedOpts         []string
			defaultOpts         map[string]interface{}
			ignoredOpts         []string
			keyPerms            map[string]string
			mandatoryOpts       []string
			actualRes           vmo.MountOpts
			err                 error
			userInput           map[string]interface{}
			mask                vmo.MountOptsMask
			validationFuncs     []vmo.UserOptsValidation
			fakeValidationFuncI *volumemountoptionsfakes.FakeUserOptsValidation
		)

		BeforeEach(func() {
			allowedOpts = []string{}
			defaultOpts = map[string]interface{}{}
			ignoredOpts = []string{}
			keyPerms = map[string]string{}
			mandatoryOpts = []string{}

			userInput = map[string]interface{}{}

			fakeValidationFuncI = &volumemountoptionsfakes.FakeUserOptsValidation{}
			validationFuncs = []vmo.UserOptsValidation{
				fakeValidationFuncI,
			}
		})

		JustBeforeEach(func() {
			mask, err = vmo.NewMountOptsMask(allowedOpts,
				defaultOpts,
				keyPerms,
				ignoredOpts,
				mandatoryOpts,
				validationFuncs...)
			Expect(err).NotTo(HaveOccurred())

			actualRes, err = vmo.NewMountOpts(userInput, mask)
		})

		Context("when given a set of options", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{
					"opt1": "val1",
					"opt2": "val2",
				}

				allowedOpts = []string{"opt1", "opt2"}
			})

			It("should return those options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualRes).To(Equal(vmo.MountOpts{
					"opt1": "val1",
					"opt2": "val2",
				}))
			})

			Context("and given a set of allowed option validations", func() {
				var (
					errorMessage1 = "errorMessage1"
					errorMessage2 = "errorMessage2"
				)

				Context("when validation check fails", func() {
					BeforeEach(func() {
						userInput = map[string]interface{}{
							"opt1": "val1",
						}

						fakeValidationFuncI.ValidateReturns(errors.New(errorMessage1))
					})

					It("should fail with a meaningful validation error", func() {
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("- validation mount options failed: %s", errorMessage1)))
					})
				})

				Context("when multiple validation checks fails", func() {
					BeforeEach(func() {
						fakeValidationFuncI.ValidateReturnsOnCall(0, errors.New(errorMessage1))
						fakeValidationFuncI.ValidateReturnsOnCall(1, errors.New(errorMessage2))
					})

					It("should fail with multiple validation errors", func() {
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("- validation mount options failed: %s, %s", errorMessage1, errorMessage2)))
					})
				})

				DescribeTable("with non string user options", func(userValue interface{}) {
					fakeValidationFuncI = &volumemountoptionsfakes.FakeUserOptsValidation{}
					userInput = map[string]interface{}{
						"opt1": userValue,
					}

					mask, err = vmo.NewMountOptsMask(allowedOpts, defaultOpts, keyPerms, ignoredOpts, mandatoryOpts, fakeValidationFuncI)
					Expect(err).NotTo(HaveOccurred())

					actualRes, err = vmo.NewMountOpts(userInput, mask)

					Expect(err).NotTo(HaveOccurred())
					expectedKey, expectedVal := fakeValidationFuncI.ValidateArgsForCall(0)
					Expect(expectedKey).To(Equal("opt1"))
					Expect(expectedVal).To(Equal(actualRes["opt1"]))
				},
					Entry("integer", 1),
					Entry("floating number", 1.0),
					Entry("null", nil),
					Entry("true", true),
					Entry("false", false),
				)

				Context("using a fake validation func", func() {
					var (
						key1, key2 string
						val1, val2 string
					)

					BeforeEach(func() {
						key1 = "key1"
						key2 = "key2"
						val1 = "val1"
						val2 = "val2"

						userInput = map[string]interface{}{
							key1: val1,
							key2: val2,
						}

						allowedOpts = []string{key1, key2}
					})

					It("should call the validation func on each user option", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(fakeValidationFuncI.ValidateCallCount()).To(Equal(2))

						key, value := fakeValidationFuncI.ValidateArgsForCall(0)
						Expect(key + value).To(Or(Equal(key1+val1), Equal(key2+val2)))

						key, value = fakeValidationFuncI.ValidateArgsForCall(1)
						Expect(key + value).To(Or(Equal(key1+val1), Equal(key2+val2)))
					})
				})

				Context("with multiple validation funcs", func() {
					var (
						key1, val1                         string
						fuzzer                             = fuzz.New().NilChance(0)
						validationFuncs1, validationFuncs2 *volumemountoptionsfakes.FakeUserOptsValidation
					)

					BeforeEach(func() {
						fuzzer.Fuzz(&key1)
						fuzzer.Fuzz(&val1)

						userInput = map[string]interface{}{
							key1: val1,
						}

						allowedOpts = []string{key1}

						validationFuncs1 = &volumemountoptionsfakes.FakeUserOptsValidation{}
						validationFuncs2 = &volumemountoptionsfakes.FakeUserOptsValidation{}
						validationFuncs = []vmo.UserOptsValidation{
							validationFuncs1,
							validationFuncs2,
						}
					})

					It("should call each validation func on each user option", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(validationFuncs1.ValidateCallCount()).To(Equal(1))
						Expect(validationFuncs2.ValidateCallCount()).To(Equal(1))

						key, value := validationFuncs1.ValidateArgsForCall(0)
						Expect(key + value).To(Equal(key1 + val1))

						key, value = validationFuncs2.ValidateArgsForCall(0)
						Expect(key + value).To(Equal(key1 + val1))
					})
				})

			})
		})

		Context("when given an empty set of opts", func() {
			It("should return those options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualRes).To(BeEmpty())
			})
		})

		Context("when some options have default values", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{
					"opt1": "val1",
					"opt2": "val2",
				}

				allowedOpts = []string{"opt1", "opt2", "opt3"}
				defaultOpts = map[string]interface{}{
					"opt2": "default2",
					"opt3": "default3",
				}
			})

			It("should return those options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualRes).To(Equal(vmo.MountOpts{
					"opt1": "val1",
					"opt2": "val2",
					"opt3": "default3",
				}))
			})

			It("should set sloppy_mount to false", func() {
				Expect(mask.SloppyMount).To(BeFalse())
			})

			It("should not mutate the mask", func() {
				m, err := vmo.NewMountOptsMask(
					[]string{"opt1", "opt2", "opt3"},
					map[string]interface{}{
						"opt2": "default2",
						"opt3": "default3",
					},
					map[string]string{},
					[]string{},
					[]string{},
					validationFuncs...,
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(mask).To(Equal(m))
			})

			Context("when there isn't any user input", func() {
				BeforeEach(func() {
					userInput = map[string]interface{}{}
				})

				It("should return the default options", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actualRes).To(Equal(vmo.MountOpts(defaultOpts)))
				})
			})
		})

		Context("when given options that are not allowed", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{
					"opt1": "val1",
					"opt2": "val2",
					"opt3": "val3",
				}

				allowedOpts = []string{"opt1"}
			})

			It("should return an error", func() {
				Expect(err.Error()).To(ContainSubstring("Not allowed options:"))
				Expect(err.Error()).To(ContainSubstring("opt2"))
				Expect(err.Error()).To(ContainSubstring("opt3"))
			})

			Context("given the sloppy_mount flag is true", func() {
				BeforeEach(func() {
					defaultOpts = map[string]interface{}{
						"sloppy_mount": "true",
					}
				})

				It("should return those options", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actualRes).To(Equal(vmo.MountOpts{
						"opt1":         "val1",
						"sloppy_mount": "true",
					}))
				})

				It("should set sloppy_mount to true", func() {
					Expect(mask.SloppyMount).To(BeTrue())
				})
			})
		})

		Context("when given a set of key permutations", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{
					"something": "some-value",
					"thing1":    "",
				}

				allowedOpts = []string{"something-else", "thing2"}
				keyPerms = map[string]string{
					"something": "something-else",
					"thing1":    "thing2",
				}
			})

			It("should create a MountOpts with the canonicalised key names", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualRes).To(Equal(vmo.MountOpts{
					"something-else": "some-value",
					"thing2":         "",
				}))
			})

			Context("when a permuted option is not allowed", func() {
				BeforeEach(func() {
					userInput = map[string]interface{}{
						"something": "some-value",
					}

					allowedOpts = []string{}
					keyPerms = map[string]string{
						"something": "something-else",
					}
				})

				It("should return an error", func() {
					Expect(err.Error()).To(Equal("- Not allowed options: something\n"))
				})
			})
		})

		Context("when given an ignored option", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{
					"something": "ignored",
				}

				ignoredOpts = []string{"something"}
			})

			It("should create a MountOpts without the ignored option", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualRes).To(BeEmpty())
			})
		})

		Context("when mandatory options are not provided", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{}

				allowedOpts = []string{"required1", "required2"}
				mandatoryOpts = []string{"required1", "required2"}
			})

			It("return an error", func() {
				Expect(err.Error()).To(ContainSubstring("Missing mandatory options: "))
				Expect(err.Error()).To(ContainSubstring("required1"))
				Expect(err.Error()).To(ContainSubstring("required2"))
			})

			Context("when given a set of key permutations", func() {
				BeforeEach(func() {
					userInput = map[string]interface{}{}

					allowedOpts = []string{"required2", "required3"}
					keyPerms = map[string]string{"required1": "required3"}
					mandatoryOpts = []string{"required2", "required3"}
				})

				It("return an error", func() {
					Expect(err.Error()).To(ContainSubstring("Missing mandatory options: "))
					Expect(err.Error()).To(ContainSubstring("required3"))
					Expect(err.Error()).To(ContainSubstring("required2"))
				})
			})
		})

		Context("when disallowed options, missing mandatory, and failed validations", func() {
			BeforeEach(func() {
				fakeValidationFuncI.ValidateReturns(errors.New("validation error"))
				allowedOpts = []string{"opt1"}
				userInput = map[string]interface{}{
					"opt1":       "val1",
					"notallowed": "foo",
				}
				mandatoryOpts = []string{"required1"}
			})

			It("returns a list of all errors", func() {
				Expect(actualRes).To(Equal(vmo.MountOpts{}))
				Expect(err).To(MatchError(
					`- validation mount options failed: validation error
- Not allowed options: notallowed
- Missing mandatory options: required1
`))
			})
		})

		Context("given int and bool options", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{
					"int":                  1,
					"int8":                 2,
					"int16":                3,
					"int32":                4,
					"int64":                5,
					"float32":              float32(1.0),
					"float64":              float64(2.0),
					"auto-traverse-mounts": true,
					"dircache":             false,
					"bool1":                true,
					"bool2":                false,
				}

				allowedOpts = []string{
					"int",
					"int8",
					"int16",
					"int32",
					"int64",
					"float32",
					"float64",
					"auto-traverse-mounts",
					"dircache",
					"bool1",
					"bool2",
				}
			})

			It("convert values to strings", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualRes).To(HaveKeyWithValue("int", "1"))
				Expect(actualRes).To(HaveKeyWithValue("int8", "2"))
				Expect(actualRes).To(HaveKeyWithValue("int16", "3"))
				Expect(actualRes).To(HaveKeyWithValue("int32", "4"))
				Expect(actualRes).To(HaveKeyWithValue("int64", "5"))
				Expect(actualRes).To(HaveKeyWithValue("float32", "1"))
				Expect(actualRes).To(HaveKeyWithValue("float64", "2"))
				Expect(actualRes).To(HaveKeyWithValue("auto-traverse-mounts", "1"))
				Expect(actualRes).To(HaveKeyWithValue("dircache", "0"))
				Expect(actualRes).To(HaveKeyWithValue("bool1", "true"))
				Expect(actualRes).To(HaveKeyWithValue("bool2", "false"))
			})
		})

		Context("given a default option that is not allowed", func() {
			BeforeEach(func() {
				userInput = map[string]interface{}{}

				allowedOpts = []string{}
				defaultOpts = map[string]interface{}{"something": "default"}
			})

			It("does not return a 'not allowed options' error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualRes).To(Equal(vmo.MountOpts{
					"something": "default",
				}))
			})
		})
	})
})
