package volume_mount_options_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVolumeMountOptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VolumeMountOptions Suite")
}
