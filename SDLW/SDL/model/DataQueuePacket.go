package m

// DataQueuePacket
/*
typedef struct SDL_DataQueuePacket
{
    size_t datalen;
	size_t startpos;
	struct SDL_DataQueuePacket *next;
	Uint8 data[SDL_VARIABLE_LENGTH_ARRAY];
} SDL_DataQueuePacket;
*/
type DataQueuePacket struct {
	Datalen  SizeT // 64-bit DLL must be used.
	Startpos SizeT // 64-bit DLL must be used.
	Next     *DataQueuePacket
	Data     [SDL_VARIABLE_LENGTH_ARRAY]Uint8 //TODO:???
}
