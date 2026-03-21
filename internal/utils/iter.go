package utils

import "iter"

func SliceSeq[T any](items []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}

// Map: T -> R
func Map[T any, R any](
	seq iter.Seq[T],
	f func(T) R,
) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Reduce: агрегирует значения
func Reduce[T any, R any](
	seq iter.Seq[T],
	init R,
	f func(R, T) R,
) R {
	result := init
	for v := range seq {
		result = f(result, v)
	}
	return result
}
