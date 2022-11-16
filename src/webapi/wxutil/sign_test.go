package wxutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckSign(t *testing.T) {
	Convey("Test CheckSign", t, func() {
		Convey("case1 failed", func() {
			caseSigns := [2]string{"ss", "s"}
			actualResp := CheckSign(caseSigns[0], caseSigns[1])
			So(actualResp, ShouldEqual, false)
		})
		Convey("case2 success", func() {
			caseSigns2 := [2]string{"ss", "ss"}
			actualResp := CheckSign(caseSigns2[0], caseSigns2[1])
			So(actualResp, ShouldEqual, true)
		})
		Convey("case3 failed", func() {
			caseSigns3 := [2]string{"", ""}
			actualResp := CheckSign(caseSigns3[0], caseSigns3[1])
			So(actualResp, ShouldEqual, false)
		})
	})
}
