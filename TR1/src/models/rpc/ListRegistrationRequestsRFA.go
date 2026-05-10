package rm

type ListRegistrationRequestsRFAParams struct {
	CommonParams
	PageRequested
}

type ListRegistrationRequestsRFAResult struct {
	CommonResult
	ItemsPaginated
}
