package rolling

type Option func(opts *rollingNumberOptions)

type rollingNumberOptions struct {
	Window int64
}

func WithRollingWindow(window int64) Option {
	return func(opts *rollingNumberOptions) {
		if window > 0 {
			opts.Window = window
		}
	}
}

func newOptions() *rollingNumberOptions {
	return &rollingNumberOptions{
		Window: defaultWindow,
	}
}

func buildOptions(opts ...Option) *rollingNumberOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}
