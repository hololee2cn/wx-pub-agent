package utils

import "testing"

func Test_StrToRGBA(t *testing.T) {
	// rgba(255, 255, 255, 255)
	rgba, err := StrToRGBA("#fffxxfff")
	if err == nil {
		t.Errorf("error str parse success, %#v", rgba)
	}

	// rgba(255, 255, 255, 255)
	rgba, err = StrToRGBA("#ffffffff")
	if err != nil {
		t.Error(err)
	}
	if rgba.R != 255 || rgba.G != 255 || rgba.B != 255 || rgba.A != 255 {
		t.Errorf("failed: %+v", rgba)
	}

	// rgba(128, 128, 128, 128)
	rgba, err = StrToRGBA("#80808080")
	if err != nil {
		t.Error(err)
	}
	if rgba.R != 128 || rgba.G != 128 || rgba.B != 128 || rgba.A != 128 {
		t.Errorf("failed: %+v", rgba)
	}

	// rgba(33, 44, 55, 66)
	rgba, err = StrToRGBA("#212C3742")
	if err != nil {
		t.Error(err)
	}
	if rgba.R != 33 || rgba.G != 44 || rgba.B != 55 || rgba.A != 66 {
		t.Errorf("failed: %+v", rgba)
	}
}
