package urlquery

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(i interface{}) (string, error) {
	return MarshalName(i, "")
}

func MarshalName(i interface{}, name string) (string, error) {
	v := reflect.ValueOf(i)
	return marshal(v, name)
}

func marshal(v reflect.Value, name string) (string, error) {
	if !v.IsValid() {
		return "", nil
	}
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			return "", nil
		}
		v = v.Elem()
		t = t.Elem()
	}
	switch {
	case isScalar(t):
		return marshalScalar(v, name)
	case isArray(t):
		return marshalArray(v, name)
	case isMap(t):
		return marshalMap(v, name)
	case isStruct(t):
		return marshalStruct(v, name)
	default:
		return "", fmt.Errorf("unknown type: %s", t.String())
	}
}

func marshalStruct(v reflect.Value, name string) (string, error) {
	t := v.Type()
	parts := []string{}
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		fv := v.Field(i)
		fn := name + fld.Name
		if tag, ok := fld.Tag.Lookup("url"); ok {
			if tag == "-" {
				continue
			}
			fn = tag
		}
		if !fv.CanSet() {
			continue
		}
		val, err := marshal(fv, fn)
		if err != nil {
			return "", err
		}
		if val != "" {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, "&"), nil
}

func marshalMap(v reflect.Value, name string) (string, error) {
	if !isString(v.Type().Key()) {
		return "", fmt.Errorf("map key must be string")
	}
	parts := []string{}
	for _, key := range v.MapKeys() {
		val, err := marshal(v.MapIndex(key), name+key.String())
		if err != nil {
			return "", err
		}
		if val != "" {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, "&"), nil
}

func marshalArray(v reflect.Value, name string) (string, error) {
	parts := []string{}
	et := v.Type().Elem()
	if isScalar(et) {
		for i := 0; i < v.Len(); i++ {
			val, _ := marshalScalar(v.Index(i), name)
			if val != "" {
				parts = append(parts, val)
			}
		}
		return strings.Join(parts, "&"), nil
	}
	for i := 0; i < v.Len(); i++ {
		newName := name + "." + strconv.Itoa(i)
		if isStruct(et) || isMap(et) {
			newName = newName + "."
		}
		val, err := marshal(v.Index(i), newName)
		if err != nil {
			return "", err
		}
		if val != "" {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, "&"), nil
}

func marshalScalar(v reflect.Value, name string) (string, error) {
	t := v.Type()
	var s string
	switch {
	case isInt(t):
		s = strconv.FormatInt(v.Int(), 10)
	case isFloat(t):
		s = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case isString(t):
		s = v.String()
	case isBool(t):
		s = strconv.FormatBool(v.Bool())
	}
	return url.QueryEscape(name) + "=" + url.QueryEscape(s), nil
}
