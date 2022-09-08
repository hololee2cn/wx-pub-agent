package repository

import "github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"

func NewSingletonRepository() {
	// repository init
	NewWXRepository(
		persistence.DefaultWxRepo(), persistence.DefaultUserRepo(), persistence.DefaultMessageRepo())
	NewAccessTokenRepository(
		persistence.DefaultAkRepo())
	NewUserRepository(
		persistence.DefaultUserRepo(), persistence.DefaultPhoneVerifyRepo())
	NewMessageRepository(
		persistence.DefaultMessageRepo(), persistence.DefaultUserRepo())
	NewTmplRepository(
		persistence.DefaultTmplRepo())
}
