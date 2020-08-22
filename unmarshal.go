package urlquery

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type values struct {
	vals url.Values
}

func (v *values) HasPrefix(s string) bool {
	for k := range v.vals {
		if strings.HasPrefix(k, s) {
			return true
		}
	}
	return false
}

func (v *values) Get(k string) string {
	ret := v.vals.Get(k)
	delete(v.vals, k)
	return ret
}

func (v *values) GetAll(k string) []string {
	ret := v.vals[k]
	delete(v.vals, k)
	return ret
}

func (v *values) Keys(pfx string) []string {
	out := []string{}
	for k := range v.vals {
		if strings.HasPrefix(k, pfx) {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func Unmarshal(i interface{}, q string) error {
	v := reflect.ValueOf(i)
	if !isPtr(v.Type()) {
		return errors.New("input should be a pointer. try with &")
	}
	vals, err := url.ParseQuery(q)
	if err != nil {
		return err
	}
	return unmarshal(v, "", &values{vals})
}

func unmarshal(v reflect.Value, name string, vals *values) error {
	t := v.Type()
	switch {
	case isScalar(t):
		return unmarshalScalar(v, vals.Get(name))
	case isArray(t):
		return unmarshalArray(v, name, vals)
	case isMap(t):
		if name != "" {
			name += "."
		}
		return unmarshalMap(v, name, vals)
	case isStruct(t):
		if name != "" {
			name += "."
		}
		return unmarshalStruct(v, name, vals)
	default:
		return fmt.Errorf("unknown type: %s", t.String())
	}
}

func unmarshalScalar(v reflect.Value, val string) error {
	t := v.Type()
	if val == "" {
		return nil
	}
	if isPtr(t) {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
	var s interface{}
	var err error
	switch {
	case isInt(t):
		s, err = strconv.ParseInt(val, 10, 64)
	case isFloat(t):
		s, err = strconv.ParseFloat(val, 64)
	case isString(t):
		s = val
	case isBool(t):
		s, err = strconv.ParseBool(val)
	}
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(s).Convert(t))
	return nil
}

func unmarshalArray(v reflect.Value, name string, vals *values) error {
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			if !vals.HasPrefix(name) {
				return nil
			}
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
	slc := reflect.MakeSlice(t, 0, 10)
	et := t.Elem()
	if isScalar(et) {
		for _, k := range vals.GetAll(name) {
			ret := reflect.New(et).Elem()
			err := unmarshalScalar(ret, k)
			if err != nil {
				return err
			}
			slc = reflect.Append(slc, ret)
		}
		v.Set(slc)
		return nil
	}
	i := 0
	for {
		newName := name + "." + strconv.Itoa(i)
		if vals.HasPrefix(newName) {
			ret := reflect.New(et).Elem()
			err := unmarshal(ret, newName, vals)
			if err != nil {
				return err
			}
			slc = reflect.Append(slc, ret)
			i++
			continue
		}
		break
	}
	v.Set(slc)
	return nil
}

func unmarshalMap(v reflect.Value, name string, vals *values) error {
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			if !vals.HasPrefix(name) {
				return nil
			}
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
	et := t.Elem()
	v.Set(reflect.MakeMap(t))
	curKN := ""
	for _, key := range vals.Keys(name) {
		kn := keyName(key, name)
		if kn == curKN {
			continue
		}
		curKN = kn
		val := reflect.New(et).Elem()
		err := unmarshal(val, name+kn, vals)
		if err != nil {
			return err
		}
		fmt.Println("map", curKN)
		v.SetMapIndex(reflect.ValueOf(kn), val)
	}
	return nil
}

func unmarshalStruct(v reflect.Value, name string, vals *values) error {
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			if !vals.HasPrefix(name) {
				return nil
			}
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		fv := v.Field(i)
		if !v.CanSet() {
			continue
		}
		fldName := fld.Name
		if tag, ok := fld.Tag.Lookup("url"); ok {
			fldName = tag
		}
		fn := name + fldName
		err := unmarshal(fv, fn, vals)
		if err != nil {
			return err
		}
	}
	return nil
}

func keyName(key, pfx string) string {
	key = strings.TrimPrefix(key, pfx)
	parts := strings.Split(key, ".")
	return parts[0]
}
