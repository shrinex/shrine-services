package slices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {
	s1 := Empty[int]()
	assert.Empty(t, s1)
}

func TestNormalize(t *testing.T) {
	in := []int{1, 2, 3}
	out := Normalize(in)

	assert.NotNil(t, out)
	assert.Equal(t, len(in), len(out))
	assert.ObjectsAreEqualValues(in, out)
}

func TestMap(t *testing.T) {
	in := []int{1, 2, 3}
	out := Map(in, func(e int) string {
		return fmt.Sprint(e)
	})

	assert.NotNil(t, out)
	assert.Equal(t, len(in), len(out))
	assert.Equal(t, "1", out[0])
	assert.Equal(t, "2", out[1])
	assert.Equal(t, "3", out[2])
}

func TestFlatMap(t *testing.T) {
	in := [][]int{{1, 1}, {2, 2, 2}, {3, 3}}
	out := FlatMap(in, func(e []int) []int {
		return e
	})

	assert.NotNil(t, out)
	assert.Equal(t, 7, len(out))
	assert.Equal(t, 1, out[0])
	assert.Equal(t, 1, out[1])
	assert.Equal(t, 2, out[2])
	assert.Equal(t, 2, out[3])
	assert.Equal(t, 2, out[4])
	assert.Equal(t, 3, out[5])
	assert.Equal(t, 3, out[6])
}

func TestFilter(t *testing.T) {
	in := []int{1, 2, 3}
	out := Filter(in, func(e int) bool {
		return e > 2
	})

	assert.Equal(t, 1, len(out))
	assert.Equal(t, 3, out[0])
}

func TestReduce(t *testing.T) {
	in := []int{1, 2, 3}
	out := Reduce(in, 0, func(r int, e int) int {
		return r + e
	})

	assert.Equal(t, 6, out)
}

func TestAsMap(t *testing.T) {
	in := []int{1, 2, 3}
	out := AsMapValuer(in, func(e int) string {
		return fmt.Sprint(e)

	}, func(e int) int {
		return e * 2
	})

	assert.Equal(t, 3, len(out))
	assert.Equal(t, 2, out["1"])
	assert.Equal(t, 4, out["2"])
	assert.Equal(t, 6, out["3"])
}

func TestAsMapConflict(t *testing.T) {
	step := 0
	in := []int{1, 2, 2, 3}
	out := AsMapValuer(in, func(e int) string {
		return fmt.Sprint(e)

	}, func(e int) int {
		step = step + 1
		return e + step
	})

	assert.Equal(t, 3, len(out))
	assert.Equal(t, 2, out["1"])
	assert.Equal(t, 5, out["2"])
	assert.Equal(t, 7, out["3"])
}

func TestGroupingBy(t *testing.T) {
	in := []int{1, 2, 2, 3}
	out := GroupingBy(in, func(e int) string {
		return fmt.Sprint(e)
	})

	assert.Equal(t, 3, len(out))
	assert.EqualValues(t, []int{1}, out["1"])
	assert.EqualValues(t, []int{2, 2}, out["2"])
	assert.EqualValues(t, []int{3}, out["3"])
}

func TestGroupingByValuer(t *testing.T) {
	in := []int{1, 2, 2, 3}
	out := GroupingByValuer(in, func(e int) string {
		return fmt.Sprint(e)
	}, func(e int) int {
		return e + e
	})

	assert.Equal(t, 3, len(out))
	assert.EqualValues(t, []int{2}, out["1"])
	assert.EqualValues(t, []int{4, 4}, out["2"])
	assert.EqualValues(t, []int{6}, out["3"])
}

func TestNestedGroupingBy(t *testing.T) {
	in := []int{1, 2, 2, 3}
	out := NestedGroupingBy(in, func(e int) string {
		return fmt.Sprint(e)
	}, func(s []int) map[string]int {
		return AsMap(s, func(e int) string {
			return fmt.Sprint(e)
		})
	})

	assert.Equal(t, 3, len(out))
	assert.EqualValues(t, map[string]int{"1": 1}, out["1"])
	assert.EqualValues(t, map[string]int{"2": 2}, out["2"])
	assert.EqualValues(t, map[string]int{"3": 3}, out["3"])
}
