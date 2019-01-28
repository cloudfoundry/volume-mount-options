package volume_mount_options_test

import (
	vmo "code.cloudfoundry.org/volume-mount-options"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeMountOptionsMask", func() {
	Describe("#NewMountOptsMask", func() {
		var (
			mask vmo.MountOptsMask

			allowedOpts   []string
			defaultOpts   map[string]string
			ignoredOpts   []string
			mandatoryOpts []string
			err           error
		)

		BeforeEach(func() {
			allowedOpts = []string{}
			defaultOpts = map[string]string{}
			ignoredOpts = []string{}
			mandatoryOpts = []string{}
		})

		JustBeforeEach(func() {
			mask, err = vmo.NewMountOptsMask(allowedOpts, defaultOpts, ignoredOpts, mandatoryOpts)
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
		})

		Context("when given a set of default options", func() {
			BeforeEach(func() {
				defaultOpts = map[string]string{
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
			BeforeEach(func() {
				defaultOpts = map[string]string{
					"sloppy_mount": "true",
				}
			})

			It("should set the SloppyMount flag", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(mask.SloppyMount).To(BeTrue())
			})

			Context("given sloppy_mount is set to an invalid value", func() {
				BeforeEach(func() {
					defaultOpts = map[string]string{
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
