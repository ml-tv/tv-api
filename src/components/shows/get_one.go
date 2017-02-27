package shows

import (
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
)

// GetOneParams represents the params needed by the GetOne handler
type GetOneParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// GetOne is an API handler to
func GetOne(req *router.Request) error {
	params := req.Params.(*GetOneParams)
	show, err := GetShow(params.ID)
	if err != nil {
		return err
	}
	if show == nil {
		return httperr.NewNotFound()
	}
	req.Ok(NewPayload(show))
	return nil
}
