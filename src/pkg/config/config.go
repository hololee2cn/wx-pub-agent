package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var mutex = &sync.RWMutex{}
var configMap = make(map[string]string)

func init() {
	Init()
}

// Init 显式Init是为了可能有需要重读配置的情况
func Init() {
	var err error
	configMap, err = parseFromINI()
	if err == nil {
		log.Info("read ini success")
		return
	}
	log.Errorf("read ini fail, err: %v,try to read env", err)
	configMap = parseFromENV()
}

// 当使用文件来配置时，等号两侧允许存在空格，字符无需使用任何转义，但开头结尾不能为空格，key，value两侧的空格都将被删除
func parseFromINI() (m map[string]string, err error) {
	a := filepath.Dir(os.Args[0])
	appPath, err := filepath.Abs(a)
	if err != nil {
		log.Error(err)
		return
	}
	confPath := filepath.Join(appPath, "etc", "app.conf")

	mutex.Lock()
	defer mutex.Unlock()
	file, err := os.Open(confPath)
	if err != nil {
		log.Println("configer open config file err:", err)
		return
	}
	defer file.Close()

	m = make(map[string]string)
	buf := bufio.NewReader(file)
	for {
		var line string
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)

		// 空行，以#开头注释，或者不是**=**格式的行都不处理
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			if err != nil {
				if err == io.EOF {
					err = nil
					break
				}
				return
			}
			continue
		}

		keyValue := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		m[key] = value
	}
	return
}

// 当使用env来配置时，key，value不可能在两侧存在空格
func parseFromENV() (m map[string]string) {
	m = make(map[string]string)
	envs := os.Environ()
	for _, v := range envs {
		keyValue := strings.SplitN(v, "=", 2)
		m[keyValue[0]] = keyValue[1]
	}
	return
}

func Int(key string) (ret int, err error) {
	value, err := getData(key)
	if err != nil {
		err = fmt.Errorf("configer get key: %v fail, err: %w", key, err)
		return
	}
	ret, err = strconv.Atoi(value)
	if err != nil {
		err = fmt.Errorf("configer get key: %v fail, raw value: %v, err: %w", key, value, err)
		return
	}
	return
}

func String(key string) (ret string, err error) {
	ret, err = getData(key)
	if err != nil {
		err = fmt.Errorf("configer get key: %v fail, err: %w", key, err)
		return
	}
	return
}

func DefaultInt(key string, defaultValue int) (ret int) {
	ret, err := Int(key)
	if err != nil {
		ret = defaultValue
	}
	return
}

func DefaultString(key string, defaultValue string) (ret string) {
	ret, err := String(key)
	if err != nil {
		ret = defaultValue
	}
	return
}

func MustString(key string) (ret string) {
	ret, err := String(key)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return
}

func getData(key string) (value string, err error) {
	mutex.RLock()
	defer mutex.RUnlock()
	value, ok := configMap[key]
	if !ok {
		err = fmt.Errorf("configer failed to get data where key = %s", key)
	}
	return
}
