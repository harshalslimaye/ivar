package jsonparser

import (
	"github.com/valyala/fastjson"
)

type JsonParser struct {
	value *fastjson.Value
}

func NewJsonParser(jp *fastjson.Value) *JsonParser {
	return &JsonParser{
		value: jp,
	}
}

func (jp *JsonParser) Exists(key string) bool {
	return jp.value.Exists(key)
}

func (jp *JsonParser) GetDependencies() map[string]string {
	return jp.GetObject("dependencies")
}

func (jp *JsonParser) GetBin() map[string]string {
	if jp.Exists("bin") {
		if jp.IsObject("bin") {
			return jp.GetObject("bin")
		}
	}
	return nil
}

func (jp *JsonParser) GetObject(key string) map[string]string {
	if !jp.Exists(key) {
		return map[string]string{}
	}

	objMap := make(map[string]string)

	values := jp.value.GetObject(key)

	values.Visit(func(k []byte, v *fastjson.Value) {
		objMap[string(k)] = string(v.GetStringBytes())
	})

	return objMap
}

func (jp *JsonParser) HasDependencies() bool {
	return jp.Exists("dependencies")
}

func (jp *JsonParser) IsObject(key string) bool {
	val := jp.value.Get(key)

	return val.Type() == fastjson.TypeObject
}

func (jp *JsonParser) IsString(key string) bool {
	val := jp.value.Get(key)

	return val.Type() == fastjson.TypeString
}

func (jp *JsonParser) GetValue(key string) string {
	return string(jp.value.GetStringBytes(key))
}

func (jp *JsonParser) GetDownloadPath() string {
	dist := jp.GetObject("dist")

	if value, exists := dist["tarball"]; exists {
		return value
	}

	return ""
}