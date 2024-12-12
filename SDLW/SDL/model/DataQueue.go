package m

// DataQueue
/*
struct SDL_DataQueue
{
    SDL_DataQueuePacket *head;
	SDL_DataQueuePacket *tail;
	SDL_DataQueuePacket *pool;
	size_t packet_size;
	size_t queued_bytes;
};
*/
// typedef struct SDL_DataQueue SDL_DataQueue;
type DataQueue struct {
	Head        *DataQueuePacket
	Tail        *DataQueuePacket
	Pool        *DataQueuePacket
	PacketSize  SizeT // 64-bit DLL must be used.
	QueuedBytes SizeT // 64-bit DLL must be used.
}
