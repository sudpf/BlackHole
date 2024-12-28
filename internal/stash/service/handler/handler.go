package handler

import (
	"BlackHole/internal/stash/service/filter"
	"BlackHole/internal/stash/service/output"
	"context"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

type MessageHandler struct {
	writers []output.Writer
	filters []filter.FilterFunc
}

func NewHandler() *MessageHandler {
	return &MessageHandler{}
}

func (mh *MessageHandler) AddWriters(writers ...output.Writer) {
	mh.writers = append(mh.writers, writers...)
}

func (mh *MessageHandler) AddFilters(filters ...filter.FilterFunc) {
	mh.filters = append(mh.filters, filters...)
}

func (mh *MessageHandler) Consume(_ context.Context, _, val string) error {
	var m map[string]interface{}
	if err := jsoniter.Unmarshal([]byte(val), &m); err != nil {
		return err
	}

	for _, filter := range mh.filters {
		if m = filter(m); m == nil {
			return nil
		}
	}

	for _, mWriter := range mh.writers {
		if err := mWriter.Write(m); err != nil {
			log.Warn("write log error:%v", err)
		}
	}

	return nil
}
