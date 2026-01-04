package copier

import (
	"fmt"
	"reflect"
)

// ShallowCopyStruct copies field values from src to dst, sharing pointers, slices, and maps.
// dst must be a pointer to a struct; src can be a struct or pointer to struct.
func ShallowCopyStruct(dst, src interface{}) error {
	if dst == nil || src == nil {
		return fmt.Errorf("copier: source and destination cannot be nil")
	}
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)
	return copyStruct(dstVal, srcVal, false)
}

// DeepCopyStruct creates an independent copy of src in dst, recursively duplicating pointers,
// slices, maps, and nested structs. dst must be a pointer to a struct; src can be a struct or pointer to struct.
func DeepCopyStruct(dst, src interface{}) error {
	if dst == nil || src == nil {
		return fmt.Errorf("copier: source and destination cannot be nil")
	}
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)
	return copyStruct(dstVal, srcVal, true)
}

func copyStruct(dst, src reflect.Value, dc bool) error {
	if dst.Kind() != reflect.Ptr || dst.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("copier: destination is not a pointer to a struct")
	}
	if src.Kind() == reflect.Ptr {
		if src.IsNil() {
			return fmt.Errorf("copier: source is nil pointer")
		}
		src = src.Elem()
	}
	if src.Kind() != reflect.Struct {
		return fmt.Errorf("copier: source is not a struct")
	}

	copyFields(dst.Elem(), src, dc)
	return nil
}

func copyFields(dst, src reflect.Value, dc bool) {
	srcType := src.Type()
	dstType := dst.Type()

	for i := 0; i < src.NumField(); i++ {
		srcTypeField := srcType.Field(i)
		if srcTypeField.Tag.Get("copier") == "-" {
			continue
		}

		srcField := src.Field(i)
		name := srcTypeField.Name

		dstTypeField, ok := dstType.FieldByName(name)
		if !ok || dstTypeField.Tag.Get("copier") == "-" {
			continue
		}

		dstField := dst.FieldByName(name)
		if dstField.IsValid() && dstField.CanSet() {
			copyValue(dstField, srcField, dc)
		}
	}
}

func copyValue(dst, src reflect.Value, dc bool) {
	if !dc {
		if src.Type().AssignableTo(dst.Type()) {
			dst.Set(src)
		} else if src.Kind() == reflect.Struct && dst.Kind() == reflect.Struct {
			copyFields(dst, src, false)
		} else if src.Type().ConvertibleTo(dst.Type()) {
			dst.Set(src.Convert(dst.Type()))
		}
		return
	}

	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
		} else {
			if dst.Kind() == reflect.Ptr {
				if dst.IsNil() {
					dst.Set(reflect.New(dst.Type().Elem()))
				}
				copyValue(dst.Elem(), src.Elem(), dc)
			} else {
				copyValue(dst, src.Elem(), dc)
			}
		}

	case reflect.Map:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
		} else {
			if dst.Kind() == reflect.Map {
				dst.Set(reflect.MakeMap(dst.Type()))
				for _, key := range src.MapKeys() {
					newKey := reflect.New(key.Type()).Elem()
					newValue := reflect.New(src.MapIndex(key).Type()).Elem()
					copyValue(newKey, key, dc)
					copyValue(newValue, src.MapIndex(key), dc)
					dst.SetMapIndex(newKey, newValue)
				}
			}
		}

	case reflect.Slice:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
		} else {
			if dst.Kind() == reflect.Slice {
				dst.Set(reflect.MakeSlice(dst.Type(), src.Len(), src.Cap()))
				for i := 0; i < src.Len(); i++ {
					copyValue(dst.Index(i), src.Index(i), dc)
				}
			}
		}

	case reflect.Struct:
		if dst.Kind() == reflect.Struct {
			copyFields(dst, src, dc)
		} else if src.Type().AssignableTo(dst.Type()) {
			dst.Set(src)
		} else if src.Type().ConvertibleTo(dst.Type()) {
			dst.Set(src.Convert(dst.Type()))
		}

	default:
		if src.Type().AssignableTo(dst.Type()) {
			dst.Set(src)
		} else if src.Type().ConvertibleTo(dst.Type()) {
			dst.Set(src.Convert(dst.Type()))
		}
	}
}
