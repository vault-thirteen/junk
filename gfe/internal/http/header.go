package http

// HTTP заголовки.
const (
	// HttpHeaderClientTimeZoneName -- HTTP заголовок :: часовой пояс клиента.
	HttpHeaderClientTimeZoneName = "X-Client-Time-Zone"

	// HttpHeaderContentType -- HTTP заголовок :: тип контента.
	HttpHeaderContentType = "Content-Type"

	// HttpHeaderAuthorization -- HTTP заголовок :: авторизация.
	HttpHeaderAuthorization = "Authorization"
)

// MIME-типы, значения для HTTP заголовков.
const (
	// HttpHeaderAuthorizationTypeBearer -- HTTP заголовок :: авторизация :: предъявитель.
	HttpHeaderAuthorizationTypeBearer = "Bearer"

	// MimeTypeApplicationJson -- MIME тип :: приложение JSON.
	MimeTypeApplicationJson = "application/json"
)
