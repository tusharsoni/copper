package clogger

// Config can be used to customize the setup for the std logger.
type Config struct {
	MinLevel Level
}

func (c Config) isValid() bool {
	return c.MinLevel > 0
}

// GetDefaultConfig provides a config for the std logger with defaults
func GetDefaultConfig() Config {
	return Config{
		MinLevel: LevelDebug,
	}
}
