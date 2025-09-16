package utils

func UnwrapJoin(err error) []error {
	u, ok := err.(interface {
		Unwrap() []error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func RecursivelyUnwrapJoin(err error) []error {
	if err == nil {
		return nil
	}

	var tmp []error

	for _, e := range UnwrapJoin(err) {
		if len(UnwrapJoin(e)) == 0 {
			tmp = append(tmp, e)
		} else {
			tmp = append(tmp, RecursivelyUnwrapJoin(e)...)
		}
	}

	if len(tmp) == 0 {
		return []error{err}
	}

	return tmp
}
