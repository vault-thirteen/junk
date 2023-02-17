package messages

const (
	MsgFSignalReceived        = "signal %v is received"
	MsgFProcessHasBeenCreated = "a process with PID=%v has been created"
	MsgFProcessHasFinished    = "a process with PID=%v has finished"
)

const (
	MsgKafkaConsumerStart      = "kafka consumer has started"
	MsgKafkaConsumerStop       = "kafka consumer has stopped"
	MsgHttpServerStarting      = "http server is starting ..."
	MsgHttpServerStopped       = "http server is stopped"
	MsgHttpServerError         = "http server error"
	MsgCriticalError           = "critical error"
	MsgTasksReceiverStart      = "tasks receiver has started"
	MsgTasksReceiverStop       = "tasks receiver has stopped"
	MsgTaskReceiverStart       = "task receiver has started"
	MsgTaskReceiverStop        = "task receiver has stopped"
	MsgServiceErrorReaderStart = "service error reader has started"
	MsgServiceErrorReaderStop  = "service error reader has stopped"
)
