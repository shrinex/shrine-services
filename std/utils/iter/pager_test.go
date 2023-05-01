package iter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPager(t *testing.T) {
	var elements []int
	for i := 0; i < 5; i++ {
		elements = append(elements, i+1)
	}

	p := NewPager(func() (int64, error) {
		return 5, nil
	}, func(offset int64, size int64) ([]int, error) {
		return elements[offset : offset+size], nil
	}, WithPageSize(1))

	idx := 0
	for {
		next, err := p.Next()
		if err == Done {
			break
		}

		idx += 1

		assert.NoError(t, err)
		assert.Equal(t, 1, len(next))
		assert.Equal(t, idx, next[0])
	}
}
