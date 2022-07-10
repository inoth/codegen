package hooks

import "go.uber.org/zap/zapcore"

func ZapHookExample(entry zapcore.Entry) error {
	go func(param zapcore.Entry) {

	}(entry)
	return nil
}
