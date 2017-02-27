package shows

import (
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
)

// DeleteParams represents the params needed by the GetOne handler
type DeleteParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// Delete is an API handler to
func Delete(req *router.Request) error {
	params := req.Params.(*DeleteParams)
	show, err := GetShow(params.ID)
	if err != nil {
		return err
	}
	if show == nil {
		return httperr.NewNotFound()
	}

	req.NoContent()
	return show.FullyDelete()
}
