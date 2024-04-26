package logger

type Config struct {
	FileName string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.
	MaxAge int

	// WithCaller is the flag to enable caller in log
	WithCaller bool

	// CallerSkip is the number of stack frames to skip to find the caller
	CallerSkip int
}
