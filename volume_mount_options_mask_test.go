package volume_mount_options_test

import (
	vmo "code.cloudfoundry.org/volume-mount-options"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeMountOptionsMask", func() {
	Describe("#NewMountOptsMask", func() {
		var (
			mask vmo.MountOptsMask

			allowedOpts   []string
			defaultOpts   map[string]interface{}
			ignoredOpts   []string
			keyPerms      map[string]string
			mandatoryOpts []string
			err           error
		)

		BeforeEach(func() {
			allowedOpts = []string{}
			defaultOpts = map[string]interface{}{}
			ignoredOpts = []string{}
			keyPerms = map[string]string{}
			mandatoryOpts = []string{}
		})

		JustBeforeEach(func() {
			mask, err = vmo.NewMountOptsMask(allowedOpts, defaultOpts, keyPerms, ignoredOpts, mandatoryOpts, nil)
		})

		Context("when given a set of allowed options", func() {
			BeforeEach(func() {
				allowedOpts = []string{"opt1", "opt2"}
			})

			It("should store those allowed options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(mask.Allowed).To(ContainElement("opt1"))
				Expect(mask.Allowed).To(ContainElement("opt2"))
			})

			It("should set the SloppyMount flag to false", func() {
				Expect(mask.SloppyMount).To(BeFalse())
			})
		})

		Context("when given a set of default options", func() {
			BeforeEach(func() {
				defaultOpts = map[string]interface{}{
					"opt2": "default2",
					"opt3": "default3",
				}
			})

			It("should store those default options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(mask.Defaults).To(HaveKeyWithValue("opt2", "default2"))
				Expect(mask.Defaults).To(HaveKeyWithValue("opt3", "default3"))
			})
		})

		Context("when given an a set of ignored option", func() {
			BeforeEach(func() {
				ignoredOpts = []string{"something"}
			})

			It("should store those ignored options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(mask.Ignored).To(ContainElement("something"))
			})
		})

		Context("when given a set of key permutations", func() {
			BeforeEach(func() {
				keyPerms = map[string]string{
					"opt1": "converted1",
					"opt2": "converted2",
				}
			})

			It("should store the key permutation information", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(mask.KeyPerms).To(Equal(keyPerms))
			})
		})

		Context("when given a set of mandatory options", func() {
			BeforeEach(func() {
				mandatoryOpts = []string{"required1", "required2"}
			})

			It("should store those mandatory options", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(mask.Mandatory).To(ContainElement("required1"))
				Expect(mask.Mandatory).To(ContainElement("required2"))
			})
		})

		Context("when given a sloppy_mount in the default options", func() {
			Context("when the string is true", func() {
				BeforeEach(func() {
					defaultOpts = map[string]interface{}{
						"sloppy_mount": "true",
					}
				})

				It("should set the SloppyMount flag", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(mask.SloppyMount).To(BeTrue())
				})
			})

			Context("when the string is false", func() {
				BeforeEach(func() {
					defaultOpts = map[string]interface{}{
						"sloppy_mount": "false",
					}
				})

				It("should not set the SloppyMount flag", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(mask.SloppyMount).To(BeFalse())
				})
			})

			Context("given sloppy_mount is set to an invalid string", func() {
				BeforeEach(func() {
					defaultOpts = map[string]interface{}{
						"sloppy_mount": "invalid",
					}
				})

				It("should return an error", func() {
					Expect(err.Error()).To(Equal(`Invalid sloppy_mount option: strconv.ParseBool: parsing "invalid": invalid syntax`))
				})
			})
		})
	})
})
