package volume_mount_options_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVolumeMountOptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VolumeMountOptions Suite")
}
