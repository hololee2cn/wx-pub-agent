// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.13
// +build go1.13

package xerrors_test

import (
	"errors"
	"testing"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/pkg/errors/internal/xerrors"
)

func TestErrorsIs(t *testing.T) {
	var errSentinel = errors.New("sentinel")

	got := errors.Is(xerrors.Errorf("%w", errSentinel), errSentinel)
	if !got {
		t.Error("got false, want true")
	}

	got = errors.Is(xerrors.Errorf("%w: %s", errSentinel, "foo"), errSentinel)
	if !got {
		t.Error("got false, want true")
	}
}
