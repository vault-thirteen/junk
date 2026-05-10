package rm

type ListUserSessionsParams struct {
	CommonParams
	PageRequested
}

type ListUserSessionsResult struct {
	CommonResult
	ItemsPaginated
}
