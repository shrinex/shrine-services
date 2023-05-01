package iter

import "github.com/pkg/errors"

// Done is returned by Pager.Next method when the iteration is
// complete; when there are no more elements to return.
var Done = errors.New("no more elements in iter")

type (
	Pager[E any] interface {
		Next() (E, error)
		EstimatedSize() int64
	}

	Option = func(options *pagerOptions)

	pagerOptions struct {
		pageSize int64
	}

	pager[E any] struct {
		err        error
		step       int64
		cursor     int64
		upperBound int64
		fetch      func(int64, int64) ([]E, error)
	}
)

func NewPager[E any](len func() (int64, error),
	fetch func(offset int64, size int64) ([]E, error),
	opts ...Option) Pager[[]E] {
	upperBound, err := len()
	options := buildOptions(opts...)
	return &pager[E]{
		err:        err,
		fetch:      fetch,
		upperBound: upperBound,
		step:       options.pageSize,
	}
}

func (p *pager[E]) EstimatedSize() int64 {
	return p.upperBound
}

func (p *pager[E]) Next() (elements []E, err error) {
	if p.err != nil {
		err = p.err
		return
	}

	if p.cursor >= p.upperBound {
		p.err = Done
		err = p.err
		return
	}

	elements, err = p.fetch(p.cursor, p.step)
	if err != nil {
		p.err = err
		err = p.err
		return
	}

	p.cursor += p.step
	return
}

func WithPageSize(pageSize int64) Option {
	return func(opts *pagerOptions) {
		opts.pageSize = pageSize
	}
}

const defaultPageSize = 256

func newOptions() *pagerOptions {
	return &pagerOptions{
		pageSize: defaultPageSize,
	}
}

func buildOptions(opts ...Option) *pagerOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	if options.pageSize <= 0 {
		options.pageSize = defaultPageSize
	}

	return options
}
