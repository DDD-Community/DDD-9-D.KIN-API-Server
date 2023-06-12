package httpx

// etc
const (
	charsetUTF8 = "charset=UTF-8"
)

// header = Content-Type
const (
	ApplicationJSON            = "application/json"
	ApplicationJSONCharsetUTF8 = ApplicationJSON + "; " + charsetUTF8
	ApplicationXML             = "application/xml"
	ApplicationXMLCharsetUTF8  = ApplicationXML + "; " + charsetUTF8
	TextXML                    = "text/xml"
	TextXMLCharsetUTF8         = TextXML + "; " + charsetUTF8
	ApplicationForm            = "application/x-www-form-urlencoded"
	MultipartForm              = "multipart/form-data"
)
