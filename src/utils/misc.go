package utils

func Ternary[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}

	return falseValue
}

func EnsureInterface(value interface{}) interface{} {
	var v interface{}

	if ptr, ok := value.(*interface{}); ok {
		v = *ptr
	} else {
		v = value
	}

	return v
}

func EnsureInterfaceSlice(value interface{}) []interface{} {
	switch value.(type) {
	case *[]interface{}:
		return *value.(*[]interface{})
	case []interface{}:
		return value.([]interface{})
	default:
		return nil
	}
}

func EnsureNonPointer[T any](value interface{}, alt T) T {

	if ptr, ok := value.(*T); ok {
		if ptr != nil {
			return *ptr
		} else {
			return alt
		}
	}

	return value.(T)
}

func InterfaceSliceToNormalSlice[T any](value []interface{}) []T {
	var items = make([]T, len(value))

	for idx, item := range value {
		items[idx], _ = item.(T)
	}

	return items
}

func InterfacePtr(value interface{}) *interface{} {
	return &value
}

func TypePointer[T any, H any](value *T, converter func(T) H) *H {
	if value == nil {
		return nil
	}

	return PtrOf(converter(*value))
}

func PtrOf[T any](value T) *T {
	return &value
}

func PtrOfCopied[T any](value T) *T {
	var copied = value
	return &copied
}

func ValueOr[T any](value *T, replacement T) T {
	if value != nil {
		return *value
	}

	return replacement
}

func SliceOrEmpty[T any](value []T) []T {
	if value != nil {
		return value
	}

	return make([]T, 0)
}

func SliceOr[T any](value []T, replacement []T) []T {
	if value != nil {
		return value
	}

	return replacement
}

func AssignValueIfNotNull[T any](value *T, target *T) {
	if value != nil && target != nil {
		*target = *value
	}
}

func MorphNullable[T any, H any](value *T, morph func(value *T) H) *H {
	if value != nil {
		return PtrOf(morph(value))
	}

	return nil
}

func ValueOrEmpty(value *string) string {
	if value != nil {
		return *value
	}

	return ""
}

func StringOrNull(value *string) *string {
	if value == nil || *value == "" {
		return nil
	}

	return value
}

func NonEmptyOr[T ~string](value T, replacement T) T {
	if value != "" {
		return value
	}

	return replacement
}

func NonEmptyPtrOr[T ~string](value *T, replacement ...*T) *T {
	if value != nil && *value != "" {
		return value
	}

	if len(replacement) > 0 {
		return replacement[0]
	} else {
		return nil
	}
}

func ToInterfaceSlice[T any](values ...T) []interface{} {
	var result = make([]interface{}, len(values))

	for i, value := range values {
		result[i] = value
	}

	return result
}

func ToAnyArray[T any](values ...T) []any {
	var result = make([]any, len(values))

	for i, value := range values {
		result[i] = value
	}

	return result
}

func ToTypedArray[T any, H any](values []T, adapter func(T) H) []H {
	var result = make([]H, len(values))

	for i, value := range values {
		result[i] = adapter(value)
	}

	return result
}

func SafeCastValue[T any](value interface{}) *T {
	if value == nil {
		return nil
	}

	if v, ok := value.(*interface{}); ok {
		if v == nil {
			return nil
		}

		value = *v
	}

	if v, ok := value.(T); ok {
		return &v
	}

	if v, ok := value.(*T); ok {
		return v
	}

	if v, ok := value.(**T); ok && v != nil {
		return *v
	}

	return value.(*T)
}
