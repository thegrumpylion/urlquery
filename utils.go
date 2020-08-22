package urlquery

import "reflect"

// go type matching

func isPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr
}

func isBool(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Bool
}

func isInt(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Int ||
		t.Kind() == reflect.Int8 ||
		t.Kind() == reflect.Int16 ||
		t.Kind() == reflect.Int32 ||
		t.Kind() == reflect.Int64 ||
		t.Kind() == reflect.Uint ||
		t.Kind() == reflect.Uint8 ||
		t.Kind() == reflect.Uint16 ||
		t.Kind() == reflect.Uint32 ||
		t.Kind() == reflect.Uint64
}

func isFloat(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Float32 ||
		t.Kind() == reflect.Float64
}

func isString(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.String
}

func isStruct(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

func isMap(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Map
}

func isArray(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Slice ||
		t.Kind() == reflect.Array
}

func isNumber(t reflect.Type) bool {
	return isInt(t) || isFloat(t)
}

func isScalar(t reflect.Type) bool {
	return isBool(t) ||
		isNumber(t) || isString(t)
}
