//go:build !windows

package webkitgtk

import (
	"github.com/ebitengine/purego"
	"github.com/pkg/errors"
)

var (
	opengl uintptr
)

func (c *defaultContext) init() error {
	lib, err := purego.Dlopen("libwebkit2gtk-4.0.so", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err == nil {
		opengl = lib
		return nil
	}

	return errors.Wrap(err, "failed to load libwebkit2gtk")
}

func (c *defaultContext) getProcAddress(name string) (uintptr, error) {
	proc, err := purego.Dlsym(opengl, name)
	if err != nil {
		return 0, errors.Wrap(err, "failed to load proc address")
	}

	return proc, nil
}
