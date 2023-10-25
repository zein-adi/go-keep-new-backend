package helpers

import (
	"errors"
)

// FindIndex mencari index dari array, mengeluarkan error bila tidak ditemukan
func FindIndex[D any](data []D, f func(D) bool) (int, error) {
	lastId := -1
	for i, datum := range data {
		if f(datum) {
			lastId = i
			break
		}
	}
	if lastId == -1 {
		return lastId, errors.New("not found")
	}
	return lastId, nil
}

// Slice mengambil array dari index skip, sejumlah take
// bila skip > data return []
// bila take mencapai akhir data, akan return sisanya
func Slice[D any](data []D, skip int, take int) []D {
	if skip >= len(data) {
		return []D{}
	}
	startIndex := skip
	endIndex := min(len(data), skip+take)
	return data[startIndex:endIndex]
}

// Reduce menjumlahkan isi (number) dari seluruh array
func Reduce[D any, E int | float32 | float64](data []D, initial E, f func(E, D) E) E {
	lastValue := initial
	for _, datum := range data {
		lastValue = f(lastValue, datum)
	}
	return lastValue
}

// Map memetakan isi array sesuai function yang dimasukkan
func Map[D any, V any](data []D, v func(D) V) []V {
	result := make([]V, 0)
	for _, datum := range data {
		result = append(result, v(datum))
	}
	return result
}

// KeyBy map array, key sesuai dengan function yang dimasukkan
func KeyBy[D any, K int | string](data []D, k func(D) K) map[K]D {
	result := make(map[K]D, len(data))
	for _, datum := range data {
		key := k(datum)
		result[key] = datum
	}
	return result
}

// KeyByMap map array dan isinya, key sesuai dengan function, dan map sesuai dengan function yang dimasukkan
func KeyByMap[D any, K int | string, V any](data []D, k func(D) K, v func(D) V) map[K]V {
	result := make(map[K]V, len(data))
	for _, datum := range data {
		key := k(datum)
		result[key] = v(datum)
	}
	return result
}

// Filter mengambil sejumlah nilai sesuai function yang dimasukkan
func Filter[D any](data []D, f func(D) bool) []D {
	result := make([]D, 0)
	for _, datum := range data {
		if f(datum) {
			result = append(result, datum)
		}
	}
	return result
}

func Unique[D int | string](data []D) []D {
	keys := make(map[D]bool, len(data))
	for _, datum := range data {
		keys[datum] = true
	}
	var result []D
	for res := range keys {
		result = append(result, res)
	}
	return result
}
