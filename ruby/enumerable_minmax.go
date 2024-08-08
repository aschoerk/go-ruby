package ruby

func (e *enumerableImpl[T]) Min(cmp Comparer[T]) (T, bool) {
	enum := e.enumerator.Clone()
	var minEl *T = nil

	for enum.HasNext() {
		tmp, _ := enum.Next()
		if minEl == nil {
			minEl = tmp
		} else {
			if cmp(*minEl, *tmp) > 0 {
				minEl = tmp
			}
		}
	}

	if minEl == nil {
		var zeroValue T
		return zeroValue, false
	} else {
		return *minEl, true
	}
}
