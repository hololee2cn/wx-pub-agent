// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"fmt"
	"os"

	"github.com/hololee2cn/captcha/internal/pkg/errors/internal/xerrors"
)

func ExampleAs() {
	_, err := os.Open("non-existing")
	if err != nil {
		var pathError *os.PathError
		if xerrors.As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		}
	}

	// Output:
	// Failed at path: non-existing
}
