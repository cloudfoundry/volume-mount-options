package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	vmo "code.cloudfoundry.org/volume-mount-options"
)

func ToKernelMountOptionString(mountOpts vmo.MountOpts) string {
	paramList := []string{}

	for k, v := range mountOpts {
		if val, err := strconv.ParseInt(v, 10, 16); err == nil {
			paramList = append(paramList, fmt.Sprintf("%s=%d", k, val))
		} else if v == "" {
			paramList = append(paramList, k)
		} else {
			paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(paramList)
	return strings.Join(paramList, ",")
}

func ParseOptionStringToMap(optionString, separator string) map[string]string {
	mountOpts := make(map[string]string, 0)

	if optionString == "" {
		return mountOpts
	}

	opts := strings.Split(optionString, ",")

	for _, opt := range opts {
		optSegments := strings.SplitN(opt, separator, 2)

		if len(optSegments) == 1 {
			mountOpts[optSegments[0]] = ""
		} else {
			mountOpts[optSegments[0]] = optSegments[1]
		}
	}

	return mountOpts
}
