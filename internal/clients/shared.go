package clients

import (
	"github.com/bytedance/sonic"
	"github.com/kasragay/backend/internal/utils"
	"github.com/valyala/fasthttp"
)

func handleError(errResp *fasthttp.Response) error {
	var errResp_ *utils.ErrorResp
	err := sonic.Unmarshal(errResp.Body(), &errResp_)
	if err != nil {
		return utils.NewInternalError(err, "failed to unmarshal error response")
	}
	return errResp_.ToError(errResp.StatusCode())
}

var _ = handleError
