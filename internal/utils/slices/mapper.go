package slices

// Map создает новый список с результатами применения функции mapper к каждому элементу из s.
func Map[A, B any](s []A, mapper func(A) B) []B {
	if s == nil {
		return nil
	}

	r := make([]B, 0, len(s))

	for _, a := range s {
		r = append(r, mapper(a))
	}

	return r
}

// MapNil создает новый список с результатами применения функции mapper к каждому элементу из s без nil.
func MapNil[A, B any](s []A, mapper func(*A) *B) []B {
	if s == nil {
		return nil
	}

	r := make([]B, 0, len(s))

	for _, a := range s {
		a := a
		if b := mapper(&a); b != nil {
			r = append(r, *b)
		}
	}

	return r
}

// MapNilError создает новый список с результатами применения функции mapper к каждому элементу из s с возвратом ошибки маппера.
func MapNilError[A, B any](s []A, mapper func(*A) (*B, error)) ([]B, error) {
	if s == nil {
		return nil, nil
	}

	r := make([]B, 0, len(s))

	for _, a := range s {
		a := a
		if b, err := mapper(&a); err != nil {
			return nil, err
		} else if b != nil {
			r = append(r, *b)
		}
	}

	return r, nil
}

// MapError создает новый список с результатами применения функции mapper к каждому элементу из s с возвратом ошибки маппера.
func MapError[A, B any](s []A, mapper func(A) (*B, error)) ([]B, error) {
	if s == nil {
		return nil, nil
	}

	r := make([]B, 0, len(s))

	for _, a := range s {
		if b, err := mapper(a); err != nil {
			return nil, err
		} else if b != nil {
			r = append(r, *b)
		}
	}

	return r, nil
}
