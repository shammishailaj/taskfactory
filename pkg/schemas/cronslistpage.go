package schemas

type CronsListPageData struct {
	AppName         string
	AppDesc         string
	RefreshInterval int64 // in milliseconds
	URL             string
}
