package response

// Response Object Interface.
type IResponseObject interface {

	// Encodes an Object in JSON Format.
	MarshalJSON() ([]byte, error)
}
