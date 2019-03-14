package config

type taskConfig struct {
	KeepDays int `toml:"keepDays"`
}

func (tc *taskConfig) FormatConfig() {
	if tc.KeepDays == 0 {
		tc.KeepDays = 7
	}
}
