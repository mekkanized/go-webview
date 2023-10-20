//go:build linux

package webkitgtk

import (
	"fmt"

	"github.com/ebitengine/purego"
)

var (
	opengl uintptr
)

func (c *defaultContext) init() error {
	lib, err := purego.Dlopen("libwebkit2gtk-4.0.so", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libwebkit2gtk: %w", err)
	}

	opengl = lib
	return nil
}

func (c *defaultContext) getProcAddress(name string) (uintptr, error) {
	proc, err := purego.Dlsym(opengl, name)
	if err != nil {
		return 0, fmt.Errorf("failed to load proc address: %w", err)
	}

	return proc, nil
}
