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
			mandatoryOpts []string
			opts          vmo.MountOpts
			err           error

			userInput map[string]interface{}
		)

		BeforeEach(func() {
			allowedOpts = []string{}
			defaultOpts = map[string]string{}
			ignoredOpts = []string{}
			mandatoryOpts = []string{}
		})

		JustBeforeEach(func() {
			var mask vmo.MountOptsMask
			mask, err = vmo.NewMountOptsMask(allowedOpts, defaultOpts, ignoredOpts, mandatoryOpts)
			Expect(err).NotTo(HaveOccurred())

			opts, err = vmo.NewMountOpts(userInput, mask)
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
				Expect(opts).To(HaveKeyWithValue("opt1", "val1"))
				Expect(opts).To(HaveKeyWithValue("opt2", "val2"))
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
				Expect(opts).To(HaveKeyWithValue("opt1", "val1"))
				Expect(opts).To(HaveKeyWithValue("opt2", "val2"))
				Expect(opts).To(HaveKeyWithValue("opt3", "default3"))
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
					Expect(opts).To(HaveKeyWithValue("opt1", "val1"))
					Expect(opts).NotTo(HaveKey("opt2"))
					Expect(opts).NotTo(HaveKey("opt3"))
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
				Expect(opts).NotTo(HaveKey("something"))
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

				allowedOpts = []string{"int", "int8", "int16", "int32", "int64", "float32", "float64", "auto-traverse-mounts", "dircache", "bool1", "bool2"}
			})

			It("convert values to strings", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(opts).To(HaveKeyWithValue("int", "1"))
				Expect(opts).To(HaveKeyWithValue("int8", "2"))
				Expect(opts).To(HaveKeyWithValue("int16", "3"))
				Expect(opts).To(HaveKeyWithValue("int32", "4"))
				Expect(opts).To(HaveKeyWithValue("int64", "5"))
				Expect(opts).To(HaveKeyWithValue("float32", "1"))
				Expect(opts).To(HaveKeyWithValue("float64", "2"))
				Expect(opts).To(HaveKeyWithValue("auto-traverse-mounts", "1"))
				Expect(opts).To(HaveKeyWithValue("dircache", "0"))
				Expect(opts).To(HaveKeyWithValue("bool1", "true"))
				Expect(opts).To(HaveKeyWithValue("bool2", "false"))
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
				Expect(opts).To(HaveKeyWithValue("something", "default"))
			})
		})

		Context("Given allowed and default params", func() {
			var (
				options      map[string]interface{}
				ignoreList   []string
				allowed      []string
				mountOptions map[string]string
				actualRes    vmo.MountOpts
				mask         vmo.MountOptsMask
			)

			BeforeEach(func() {
				options = make(map[string]interface{}, 0)
				ignoreList = make([]string, 0)

				allowed = []string{"sloppy_mount", "nfs_uid", "nfs_gid", "allow_other", "uid", "gid", "auto-traverse-mounts", "dircache", "foo", "bar", "flo"}
				mountOptions = map[string]string{
					"nfs_uid": "1003",
					"nfs_gid": "1001",
					"uid":     "1004",
					"gid":     "1002",
				}

				mask, err = vmo.NewMountOptsMask(allowed, mountOptions, ignoreList, []string{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("should flow sloppy_mount as disabled", func() {
				Expect(mask.SloppyMount).To(BeFalse())
			})

			Context("Given empty arbitrary params and share without any params", func() {
				BeforeEach(func() {
					actualRes, err = vmo.NewMountOpts(options, mask)
				})

				It("should return nil result on setting end users'entries", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should pass the default options into the MountOptions struct", func() {
					for k, exp := range mountOptions {
						Expect(inMapInt(actualRes, k, exp)).To(BeTrue())
					}

					for k, exp := range actualRes {
						Expect(inMapInt(mountOptions, k, exp)).To(BeTrue())
					}
				})
			})
		})
	})
})

func MatchExactly(actual, expected map[string]string) bool {
	for k, exp := range expected {
		Expect(inMapInt(actual, k, exp)).To(BeTrue())
	}

	for k, exp := range actual {
		Expect(inMapInt(expected, k, exp)).To(BeTrue())
	}
}

func inMapInt(list map[string]string, key string, val string) bool {
	for k, v := range list {
		if k != key {
			continue
		}

		if v == val {
			return true
		} else {
			return false
		}
	}

	return false
}
