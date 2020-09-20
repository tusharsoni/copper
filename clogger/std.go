package clogger

import (
	"errors"
	"log"

	"go.uber.org/fx"

	"github.com/tusharsoni/copper/cerror"
)

type NewParams struct {
	fx.In

	Config Config `optional:"true"`
}

func New(p NewParams) Logger {
	if !p.Config.isValid() {
		p.Config = GetDefaultConfig()
	}

	return &stdLogger{
		tags:   make(map[string]interface{}),
		config: p.Config,
	}
}

type stdLogger struct {
	tags   map[string]interface{}
	config Config
}

func (s *stdLogger) WithTags(tags map[string]interface{}) Logger {
	return &stdLogger{
		tags:   mergeTags(s.tags, tags),
		config: s.config,
	}
}

func (s *stdLogger) Debug(msg string) {
	s.log(LevelDebug, errors.New(msg))
}

func (s *stdLogger) Info(msg string) {
	s.log(LevelInfo, errors.New(msg))
}

func (s *stdLogger) Warn(msg string, err error) {
	s.log(LevelWarn, cerror.New(err, msg, nil))
}

func (s *stdLogger) Error(msg string, err error) {
	s.log(LevelError, cerror.New(err, msg, nil))
}

func (s *stdLogger) log(lvl Level, err error) {
	if lvl < s.config.MinLevel {
		return
	}

	log.Printf("[%s] %s", lvl.String(), cerror.WithTags(err, s.tags).Error())
}
