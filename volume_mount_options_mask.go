package volume_mount_options

import (
	"strconv"

	werror "github.com/pkg/errors"
)

type MountOptsMask struct {
	// set of options that are allowed to be provided by the user
	Allowed []string

	// set of default values that will be used if not otherwise provided
	Defaults map[string]string

	// set of options that, if provided,  will be silently ignored
	Ignored []string

	// set of options that must be provided
	Mandatory []string

	SloppyMount bool
}

func NewMountOptsMask(allowed []string, defaults map[string]string, ignored, mandatory []string) (MountOptsMask, error) {
	mask := MountOptsMask{
		Allowed:   allowed,
		Defaults:  defaults,
		Ignored:   ignored,
		Mandatory: mandatory,
	}

	if defaults == nil {
		mask.Defaults = make(map[string]string)
	}

	if v, ok := defaults["sloppy_mount"]; ok {
		var err error
		mask.SloppyMount, err = strconv.ParseBool(v)

		if err != nil {
			return MountOptsMask{}, werror.Wrap(err, "Invalid sloppy_mount option")
		}
	}

	return mask, nil
}
