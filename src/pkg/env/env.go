package env

import (
	"flag"
	"fmt"
	"strings"
)

var (
	active Environment
	dev    Environment = &environment{value: "dev"}
	fat    Environment = &environment{value: "fat"}
	uat    Environment = &environment{value: "uat"}
	pro    Environment = &environment{value: "pro"}
)

type Environment interface {
	Value() string
	IsDev() bool
	IsFat() bool
	IsUat() bool
	IsPro() bool
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsDev() bool {
	return e.value == dev.Value()
}

func (e *environment) IsFat() bool {
	return e.value == fat.Value()
}

func (e *environment) IsUat() bool {
	return e.value == uat.Value()
}

func (e *environment) IsPro() bool {
	return e.value == pro.Value()
}

func init() {
	env := flag.String("env", "", "请输入运行环境:\n dev:开发环境\n fat:测试环境\n uat:预上线环境\n pro:正式环境\n 系统默认配置为开发环境，其他环境配置需要自行配置")
	flag.Parse()

	switch strings.ToLower(strings.TrimSpace(*env)) {
	case dev.Value():
		active = dev
	case fat.Value():
		active = fat
	case uat.Value():
		active = uat
	case pro.Value():
		active = pro
	default:
		active = dev
		fmt.Println("Warning: '-env' cannot be found, or it is illegal. The default 'dev' will be used.")
	}
}

// Active 当前配置的env
func Active() Environment {
	return active
}
