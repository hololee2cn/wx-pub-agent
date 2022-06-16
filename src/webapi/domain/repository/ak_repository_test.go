package repository

import (
	"context"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
	"testing"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccessTokenRepository_GetAccessToken(t *testing.T) {
	Convey("Test AccessTokenRepository_GetAccessToken", t, func() {
		Convey("case1 get ak from redis", func() {
			persistence.NewAkRepo()
			NewAccessTokenRepository(persistence.DefaultAkRepo())
			a := DefaultAccessTokenRepository()
			monkey.Patch((*persistence.AkRepo).GetAccessTokenFromRedis, func(_ *persistence.AkRepo, ctx context.Context) (string, error) {
				return "test", nil
			})
			defer monkey.Unpatch((*persistence.AkRepo).GetAccessTokenFromRedis)
			mockCtx := context.Background()
			actualResp, err := a.GetAccessToken(mockCtx)
			expectedResp := "test"
			So(actualResp, ShouldEqual, expectedResp)
			So(err, ShouldBeNil)
		})
		Convey("case2 get ak from remote api", func() {
			persistence.NewAkRepo()
			NewAccessTokenRepository(persistence.DefaultAkRepo())
			a := DefaultAccessTokenRepository()
			monkey.Patch((*persistence.AkRepo).GetAccessTokenFromRedis, func(_ *persistence.AkRepo, ctx context.Context) (string, error) {
				return "", nil
			})
			defer monkey.Unpatch((*persistence.AkRepo).GetAccessTokenFromRedis)
			monkey.Patch((*AccessTokenRepository).FreshAccessToken, func(_ *AccessTokenRepository, ctx context.Context) (string, error) {
				return "test", nil
			})
			defer monkey.Unpatch((*AccessTokenRepository).FreshAccessToken)
			mockCtx := context.Background()
			actualResp, err := a.GetAccessToken(mockCtx)
			expectedResp := "test"
			So(actualResp, ShouldEqual, expectedResp)
			So(err, ShouldBeNil)
		})
	})
}
