package pool

// Options is used to pass option in when create pool.
type Options struct {
	eachCPU float64 // CPU resources each worker use.
	eachMem uint64  // Amount of memory each worker use in KiB.
	maxCPU  float64 // Maximum CPU this pool can allocate.
	maxMem  uint64  // Maximum amount of memory this pool can allocate in KiB.
}

// Option is a option to create pool.
type Option func(*Options)

// WithEachCPU sets the CPU resources allocated to each worker created by pool.
func WithEachCPU(val float64) Option {
	return func(args *Options) {
		args.eachCPU = val
	}
}

// WithEachMem sets the amount of memory allocated to each worker created by pool.
func WithEachMem(val uint64) Option {
	return func(args *Options) {
		args.eachMem = val
	}
}

// WithMaxCPU sets the maximum CPU resources pool can allocate
func WithMaxCPU(val float64) Option {
	return func(args *Options) {
		args.maxCPU = val
	}
}

// WithMaxMem sets the maximum amount of memory pool can allocate
func WithMaxMem(val uint64) Option {
	return func(args *Options) {
		args.maxMem = val
	}
}
