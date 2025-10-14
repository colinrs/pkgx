package asynqueue

import (
	"strings"

	"github.com/rs/xid"
)

func getTaskId(eventName, traceID string) string {
	if traceID == "" {
		traceID = xid.New().String()
	}
	return strings.Join([]string{eventName, traceID, xid.New().String()}, taskIdDelimiter)
}

func getTraceIDFromTaskId(taskId string) string {
	ids := strings.Split(taskId, taskIdDelimiter)
	if len(ids) == 3 {
		_, traceID, _ := ids[0], ids[1], ids[2]
		return traceID
	}
	return ""
}
