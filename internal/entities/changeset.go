package entities

type OptField[T any] struct {
	value T
	isSet bool
}

func (f *OptField[T]) Set(value T) {
	f.value = value
	f.isSet = true
}

func (f *OptField[T]) Unset() {
	f.isSet = false
}

func (f *OptField[T]) Get() (T, bool) {
	return f.value, f.isSet
}
