package volume_mount_options_test

import (
	vmo "code.cloudfoundry.org/volume-mount-options"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeMountOptions", func() {
	Describe("#NewMountOpts", func() {
		var (
			allowedOpts   []string
			defaultOpts   map[string]string
			ignoredOpts   []string
			keyPerms      map[string]string
			mandatoryOpts []string
			actualRes     vmo.MountOpts
			err           error

			userInput map[string]interface{}
			mask      vmo.MountOptsMask
		)

		BeforeEach(func() {
			allowedOpts = []string{}
			defaultOpts = map[string]string{}
			ignoredOpts = []string{}
			keyPerms = map[string]string{}
			mandatoryOpts = []string{}

			userInput = map[string]interface{}{}
		})

		JustBeforeEach(func() {
			mask, err = vmo.NewMountOptsMask(allowedOpts, defaultOpts, keyPerms, ignoredOpts, mandatoryOpts)
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
				Expect(map[string]string(actualRes)).To(Equal(map[string]string{
					"opt1": "val1",
					"opt2": "val2",
				}))
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
				defaultOpts = map[string]string{
					"opt2": "default2",
					"opt3": "default3",
				}
			})

			It("should return those options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(map[string]string(actualRes)).To(Equal(map[string]string{
					"opt1": "val1",
					"opt2": "val2",
					"opt3": "default3",
				}))
			})

			It("should set sloppy_mount to false", func() {
				Expect(mask.SloppyMount).To(BeFalse())
			})

			Context("when there isnt any user input", func() {
				BeforeEach(func() {
					userInput = map[string]interface{}{}
				})
				It("should return the default options", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(map[string]string(actualRes)).To(Equal(defaultOpts))
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
					defaultOpts = map[string]string{
						"sloppy_mount": "true",
					}
				})

				It("should return those options", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(map[string]string(actualRes)).To(Equal(map[string]string{
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
				Expect(map[string]string(actualRes)).To(Equal(map[string]string{
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
					Expect(err.Error()).To(Equal("Not allowed options: something"))
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
				Expect(map[string]string(actualRes)).To(BeEmpty())
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
				defaultOpts = map[string]string{"something": "default"}
			})

			It("does not return a 'not allowed options' error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(map[string]string(actualRes)).To(Equal(map[string]string{
					"something": "default",
				}))
			})
		})
	})
})
