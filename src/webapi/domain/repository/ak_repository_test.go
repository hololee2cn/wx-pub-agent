package repository

import (
	"context"
	"testing"

	monkey "github.com/agiledragon/gomonkey/v2"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccessTokenRepository_GetAccessToken(t *testing.T) {
	var (
		s             *persistence.AkRepo
		ak            *AccessTokenRepository
		defaultAkRepo *persistence.AkRepo
	)
	defaultAkRepo = new(persistence.AkRepo)
	defaultAkRepo.Redis = nil
	NewAccessTokenRepository(defaultAkRepo)
	a := DefaultAccessTokenRepository()
	Convey("Test AccessTokenRepository_GetAccessToken", t, func() {
		Convey("case1 get ak from redis", func() {
			patches := monkey.ApplyMethodFunc(s, "GetAccessTokenFromRedis", func(_ context.Context) (string, error) {
				return "test", nil
			})
			defer patches.Reset()
			mockCtx := context.Background()
			actualResp, err := a.GetAccessToken(mockCtx)
			expectedResp := "test"
			So(actualResp, ShouldEqual, expectedResp)
			So(err, ShouldBeNil)
		})
		Convey("case2 get ak from remote api", func() {
			patches := monkey.ApplyMethodFunc(s, "GetAccessTokenFromRedis", func(ctx context.Context) (string, error) {
				return "", nil
			})
			defer patches.Reset()
			patches = monkey.ApplyMethodFunc(ak, "FreshAccessToken", func(ctx context.Context) (string, error) {
				return "test", nil
			})
			defer patches.Reset()
			mockCtx := context.Background()
			actualResp, err := a.GetAccessToken(mockCtx)
			expectedResp := "test"
			So(actualResp, ShouldEqual, expectedResp)
			So(err, ShouldBeNil)
		})
	})
}
