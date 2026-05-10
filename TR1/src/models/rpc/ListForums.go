package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ListForumsParams struct {
	CommonParams
}

type ListForumsResult struct {
	CommonResult
	Forums []cm.Forum `json:"forums"`
}
