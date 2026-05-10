package rm

type ListUsersParams struct {
	CommonParams
	PageRequested
}

type ListUsersResult struct {
	CommonResult
	ItemsPaginated
}
