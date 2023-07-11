package webkitgtk

import "github.com/pkg/errors"

type procAddressGetter struct {
	ctx *defaultContext
	err error
}

func (p *procAddressGetter) get(name string) uintptr {
	if p.err != nil {
		return 0
	}

	proc, err := p.ctx.getProcAddress(name)
	if err != nil {
		p.err = errors.Wrapf(err, "webkit2gtk: %s is missing", name)
		return 0
	}
	if proc == 0 {
		p.err = errors.Errorf("webkit2gtk: %s is missing", name)
		return 0
	}

	return proc
}
