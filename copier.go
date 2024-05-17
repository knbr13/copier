package copier

import (
	"fmt"
	"reflect"
)

const (
	ErrNotAPointerToStructDestination = "not a pointer to a struct destination"
	ErrNotAStructSource               = "not a struct source"
)

func copyStruct(dst, src reflect.Value, dc bool) error {
	if dst.Kind() != reflect.Ptr || dst.Elem().Kind() != reflect.Struct {
		return fmt.Errorf(ErrNotAPointerToStructDestination)
	}
	if src.Kind() != reflect.Struct {
		return fmt.Errorf(ErrNotAStructSource)
	}

	dstElm := dst.Elem()
	dt := dstElm.Type()

	for i := 0; i < dst.Elem().NumField(); i++ {
		df := dstElm.Field(i)
		sf := src.FieldByName(dt.Field(i).Name)

		if !df.CanSet() || df.Kind() != sf.Kind() {
			continue
		}

		copyValue(df, sf, dc)
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

	case reflect.Chan:
		dst.Set(reflect.MakeChan(src.Type(), src.Cap()))

	default:
		dst.Set(src)
	}
}
