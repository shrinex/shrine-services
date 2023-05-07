// Package slices provides function to manipulate slice
package slices

// Empty returns an empty slice of element type E
func Empty[E any]() []E {
	return []E{}
}

// IsEmpty returns true if s contains no element
func IsEmpty[E any](s []E) bool {
	return len(s) == 0
}

// IsNotEmpty returns true if s contains element
func IsNotEmpty[E any](s []E) bool {
	return len(s) > 0
}

// Normalize returns a slice as input, but with element type changed
func Normalize[E any](s []E) []any {
	if len(s) == 0 {
		return Empty[any]()
	}

	result := make([]any, 0, len(s))
	for _, e := range s {
		result = append(result, e)
	}

	return result
}

// ForEach iterate over the input slice’s elements
func ForEach[E any](s []E, fn func(E, *bool)) {
	stop := false
	for _, e := range s {
		fn(e, &stop)
		if stop {
			break
		}
	}
}

// MutateEach iterate over the input slice’s elements
func MutateEach[E any](s []E, fn func(*E, *bool)) {
	stop := false
	for i := 0; i < len(s); i++ {
		fn(&s[i], &stop)
		if stop {
			break
		}
	}
}

// Map returns a slice containing the results of mapping the given
// closure over the input slice’s elements
func Map[E, T any](s []E, transform func(E) T) []T {
	if len(s) == 0 {
		return Empty[T]()
	}

	result := make([]T, 0, len(s))
	for _, e := range s {
		result = append(result, transform(e))
	}

	return result
}

// FlatMap returns a slice containing the concatenated results of calling
// the given transformation with each element of the input slice
func FlatMap[E, T any](s []E, transform func(E) []T) []T {
	if len(s) == 0 {
		return Empty[T]()
	}

	result := make([]T, 0, len(s))
	for _, e := range s {
		result = append(result, transform(e)...)
	}

	return result
}

// Filter generates a slice that contains the elements that satisfy a predicate
func Filter[E any](s []E, predicate func(E) bool) []E {
	if len(s) == 0 {
		return Empty[E]()
	}

	var result []E
	for _, e := range s {
		if predicate(e) {
			result = append(result, e)
		}
	}

	return result
}

// Reduce returns the result of combining the elements of the slice using the given closure
func Reduce[E, T any](s []E, initialResult T, nextPartialResult func(T, E) T) T {
	if len(s) == 0 {
		return initialResult
	}

	var result = initialResult
	for _, e := range s {
		result = nextPartialResult(result, e)
	}

	return result
}

// AsMap transforms slice into map
func AsMap[E any, K comparable](s []E, keyer func(E) K) map[K]E {
	return AsMapValuer(s, keyer, Identity[E])
}

func AsMapMerger[E any, K comparable](s []E, keyer func(E) K, merger func(E, E) E) map[K]E {
	return AsMapValuerMerger(s, keyer, Identity[E], merger)
}

// AsMapValuer transforms slice into map
func AsMapValuer[E, V any, K comparable](s []E, keyer func(E) K, valuer func(E) V) map[K]V {
	return AsMapValuerMerger(s, keyer, valuer, func(lhs V, rhs V) V {
		return rhs
	})
}

// AsMapValuerMerger transforms slice into map
func AsMapValuerMerger[E, V any, K comparable](s []E, keyer func(E) K, valuer func(E) V, merger func(V, V) V) map[K]V {
	return Collect[E, map[K]V](s, make(map[K]V), func(m *map[K]V, e E) {
		if v, ok := (*m)[keyer(e)]; ok {
			(*m)[keyer(e)] = merger(v, valuer(e))
		} else {
			(*m)[keyer(e)] = valuer(e)
		}
	})
}

// AsSet transforms slice into set
func AsSet[E any, K comparable](s []E, keyer func(E) K) map[K]struct{} {
	return AsMapValuer(s, keyer, func(e E) struct{} {
		return struct{}{}
	})
}

// Collect returns the result of combining the elements of the slice using the given closure
func Collect[E, T any](s []E, initialResult T, updateAccumulatingResult func(*T, E)) T {
	if len(s) == 0 {
		return initialResult
	}

	var result = initialResult
	for _, e := range s {
		updateAccumulatingResult(&result, e)
	}

	return result
}

// AnyMatch returns whether any elements of the slice match the provided predicate
func AnyMatch[E any](s []E, predicate func(E) bool) bool {
	if len(s) == 0 {
		return false
	}

	for _, e := range s {
		if predicate(e) {
			return true
		}
	}

	return false
}

// AllMatch returns whether all elements of the slice match the provided predicate
func AllMatch[E any](s []E, predicate func(E) bool) bool {
	if len(s) == 0 {
		return false
	}

	for _, e := range s {
		if !predicate(e) {
			return false
		}
	}

	return true
}

// NonMatch returns whether no elements of the slice match the provided predicate
func NonMatch[E any](s []E, predicate func(E) bool) bool {
	if len(s) == 0 {
		return true
	}

	for _, e := range s {
		if predicate(e) {
			return false
		}
	}

	return true
}

// Min returns the minimum element of the slice according to the provided comparator
func Min[E any](s []E, comparator func(E, E) int) (result E, ok bool) {
	if len(s) == 0 {
		return
	}

	ok = true
	result = s[0]
	for _, e := range s {
		if comparator(e, result) < 0 {
			result = e
		}
	}

	return
}

// Max returns the maximum element of the slice according to the provided comparator
func Max[E any](s []E, comparator func(E, E) int) (result E, ok bool) {
	if len(s) == 0 {
		return
	}

	ok = true
	result = s[0]
	for _, e := range s {
		if comparator(e, result) > 0 {
			result = e
		}
	}

	return
}

// GroupingBy implementing a "group by" operation on input elements of type E
// grouping elements according to a classification function, and returning the results in a map
func GroupingBy[E any, K comparable](s []E, classifier func(E) K) map[K][]E {
	return GroupingByValuer(s, classifier, Identity[E])
}

// GroupingByValuer implementing a "group by" operation on input elements of type E
// grouping elements according to a classification function, and returning the results in a map
func GroupingByValuer[E, V any, K comparable](s []E, classifier func(E) K, valuer func(E) V) map[K][]V {
	result := make(map[K][]V)
	for _, e := range s {
		key := classifier(e)
		slice, ok := result[key]
		if !ok {
			slice = make([]V, 0)
		}
		slice = append(slice, valuer(e))
		result[key] = slice
	}

	return result
}

// NestedGroupingBy implementing a "group by" operation on input elements of type E
// grouping elements according to a classification function, and returning the results in a map
func NestedGroupingBy[E, V any, K comparable](s []E, classifier func(E) K, valuer func([]E) V) map[K]V {
	if len(s) == 0 {
		return make(map[K]V)
	}

	group := make(map[K][]E)
	for _, e := range s {
		key := classifier(e)
		slice, ok := group[key]
		if !ok {
			slice = make([]E, 0)
		}
		slice = append(slice, e)
		group[key] = slice
	}

	result := make(map[K]V)
	for key, value := range group {
		result[key] = valuer(value)
	}

	return result
}

func Identity[E any](e E) E {
	return e
}
