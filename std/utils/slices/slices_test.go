package slices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {
	s1 := Empty[int]()
	assert.Empty(t, s1)
	assert.True(t, IsEmpty(s1))
	assert.False(t, IsNotEmpty(s1))
}

func TestNormalize(t *testing.T) {
	in := []int{1, 2, 3}
	out := Normalize(in)

	assert.NotNil(t, out)
	assert.Equal(t, len(in), len(out))
	assert.ObjectsAreEqualValues(in, out)
}

func TestForEach(t *testing.T) {
	in := []int{1, 2, 3}
	var out []int
	ForEach(in, func(i int, stop *bool) {
		i = i + 1
		out = append(out, i)
	})

	assert.EqualValues(t, []int{1, 2, 3}, in)
	assert.EqualValues(t, []int{2, 3, 4}, out)
}

func TestMutateEach(t *testing.T) {
	in := []int{1, 2, 3}
	MutateEach(in, func(i *int, stop *bool) {
		*i = *i + 1
	})

	assert.EqualValues(t, []int{2, 3, 4}, in)
}

func TestMutateEachWithStop(t *testing.T) {
	in := []int{1, 2, 3}
	MutateEach(in, func(i *int, stop *bool) {
		if *i == 2 {
			*stop = true
			return
		}
		*i = *i + 1
	})

	assert.EqualValues(t, []int{2, 2, 3}, in)
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
	assert.Equal(t, 5, out["2"]) // last one wins
	assert.Equal(t, 7, out["3"])
}

func TestAsMapValuerMerger(t *testing.T) {
	step := 0
	in := []int{1, 2, 2, 3}
	out := AsMapValuerMerger(in, func(e int) string {
		return fmt.Sprint(e)
	}, func(e int) int {
		step = step + 1
		return e + step
	}, func(lhs int, rhs int) int {
		return lhs
	})

	assert.Equal(t, 3, len(out))
	assert.Equal(t, 2, out["1"])
	assert.Equal(t, 4, out["2"]) // first one wins
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

func TestAnyMatch(t *testing.T) {
	in := []int{1, 2, 2, 3}

	assert.True(t, AnyMatch(in, func(e int) bool {
		return e == 3
	}))
	assert.False(t, AnyMatch(in, func(e int) bool {
		return e == 4
	}))
}

func TestAllMatch(t *testing.T) {
	in := []int{1, 2, 2, 3}

	assert.True(t, AllMatch(in, func(e int) bool {
		return e > 0
	}))
	assert.True(t, AllMatch(in, func(e int) bool {
		return e < 4
	}))
	assert.False(t, AllMatch(in, func(e int) bool {
		return e > 1
	}))
}

func TestNonMatch(t *testing.T) {
	in := []int{1, 2, 2, 3}

	assert.True(t, NonMatch(in, func(e int) bool {
		return e < 0
	}))
	assert.True(t, NonMatch(in, func(e int) bool {
		return e > 4
	}))
	assert.False(t, NonMatch(in, func(e int) bool {
		return e > 1
	}))
}

func TestMinWithEmptySlice(t *testing.T) {
	var in []int

	min, ok := Min(in, func(lhs int, rhs int) int {
		return lhs - rhs
	})
	assert.False(t, ok)
	assert.Equal(t, 0, min)
}

func TestMin(t *testing.T) {
	in := []int{1, 2, 3}

	min, ok := Min(in, func(lhs int, rhs int) int {
		return lhs - rhs
	})
	assert.True(t, ok)
	assert.Equal(t, 1, min)
}

func TestMaxWithEmptySlice(t *testing.T) {
	var in []int

	nax, ok := Max(in, func(lhs int, rhs int) int {
		return lhs - rhs
	})
	assert.False(t, ok)
	assert.Equal(t, 0, nax)
}

func TestMax(t *testing.T) {
	in := []int{1, 2, 3}

	max, ok := Max(in, func(lhs int, rhs int) int {
		return lhs - rhs
	})
	assert.True(t, ok)
	assert.Equal(t, 3, max)
}
