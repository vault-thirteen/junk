package m

type LogPriority = Int

const (
	LOG_PRIORITY_VERBOSE  = LogPriority(1)
	LOG_PRIORITY_DEBUG    = LogPriority(2)
	LOG_PRIORITY_INFO     = LogPriority(3)
	LOG_PRIORITY_WARN     = LogPriority(4)
	LOG_PRIORITY_ERROR    = LogPriority(5)
	LOG_PRIORITY_CRITICAL = LogPriority(6)
	NUM_LOG_PRIORITIES    = LogPriority(7)
)
