package copier

import (
	"fmt"
	"reflect"
)

func ShallowCopyStruct(dst, src interface{}) error {
	if dst == nil || src == nil {
		return fmt.Errorf("source and destination cannot be nil")
	}
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)
	return copyStruct(dstVal, srcVal, false)
}

func DeepCopyStruct(dst, src interface{}) error {
	if dst == nil || src == nil {
		return fmt.Errorf("source and destination cannot be nil")
	}
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)
	return copyStruct(dstVal, srcVal, true)
}

func copyStruct(dst, src reflect.Value, dc bool) error {
	if dst.Kind() != reflect.Ptr || dst.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination is not a pointer to a struct")
	}
	if src.Kind() != reflect.Struct {
		return fmt.Errorf("source is not a struct")
	}

	dstElem := dst.Elem()
	dstElemType := dstElem.Type()

	for i := 0; i < dstElem.NumField(); i++ {
		dstField := dstElem.Field(i)
		srcField := src.FieldByName(dstElemType.Field(i).Name)

		if !dstField.IsValid() || !dstField.CanSet() || srcField.Kind() != dstField.Kind() {
			continue
		}

		copyValue(dstField, srcField, dc)
	}

	return nil
}

func copyValue(dst, src reflect.Value, dc bool) {
	if !dc {
		dst.Set(src)
		return
	}

	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
		} else {
			dst.Set(reflect.New(src.Elem().Type()))
			copyValue(dst.Elem(), src.Elem(), dc)
		}

	case reflect.Map:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
		} else {
			dst.Set(reflect.MakeMap(src.Type()))
			for _, key := range src.MapKeys() {
				newKey := reflect.New(key.Type()).Elem()
				newValue := reflect.New(src.MapIndex(key).Type()).Elem()
				copyValue(newKey, key, dc)
				copyValue(newValue, src.MapIndex(key), dc)
				dst.SetMapIndex(newKey, newValue)
			}
		}

	case reflect.Slice:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
		} else {
			dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
			for i := 0; i < src.Len(); i++ {
				copyValue(dst.Index(i), src.Index(i), dc)
			}
		}

	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			copyValue(dst.Field(i), src.Field(i), dc)
		}

	default:
		dst.Set(src)
	}
}
