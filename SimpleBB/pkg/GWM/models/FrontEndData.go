package models

// FrontEndData is auxiliary data for front end.
type FrontEndData struct {
	AdminHtmlPage FrontEndFileData
	AdminJs       FrontEndFileData
	ApiJs         FrontEndFileData
	ArgonJs       FrontEndFileData
	ArgonWasm     FrontEndFileData
	BppJs         FrontEndFileData
	CssStyles     FrontEndFileData
	FavIcon       FrontEndFileData
	IndexHtmlPage FrontEndFileData
	LoaderScript  FrontEndFileData
}
