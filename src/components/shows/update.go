package shows

import (
	"time"

	"strings"

	"fmt"

	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
)

// UpdateParams represents the params needed by the Update handler
type UpdateParams struct {
	ID            string  `from:"url" json:"id" params:"uuid"`
	Name          string  `from:"form" json:"name" params:"trim"`
	OriginalName  string  `from:"form" json:"original_name" params:"trim"`
	Synopsis      string  `from:"form" json:"synopsis" params:"trim"`
	Status        *int    `from:"form" json:"status"`
	DayOfWeek     *int    `from:"form" json:"day_of_week"`
	ReturningDate *string `from:"form" json:"returning_date"`
	Website       *string `from:"form" json:"website" params:"trim"`
	Wikipedia     *string `from:"form" json:"wikipedia" params:"trim"`
	ExtraLink     *string `from:"form" json:"extra_link" params:"trim"`
	OnNetflix     *bool   `from:"form" json:"on_netflix" params:"trim"`
}

// Update is an API handler to
func Update(req *router.Request) error {
	params := req.Params.(*UpdateParams)
	show, err := GetShow(params.ID)
	if err != nil {
		return err
	}
	if show == nil {
		return httperr.NewNotFound()
	}

	if params.Name != "" {
		show.Name = params.Name
	}

	if params.OriginalName != "" {
		show.OriginalName = params.OriginalName
	}

	if params.Synopsis != "" {
		show.Synopsis = params.Synopsis
	}

	if params.Status != nil {
		if *params.Status < 0 || *params.Status >= ShowStatusEndOfList {
			return httperr.NewBadRequest("invalid status: %d", *params.Status)
		}
		show.Status = *params.Status
	}

	if params.DayOfWeek != nil {
		if *params.DayOfWeek < 0 || *params.DayOfWeek > 6 {
			return httperr.NewBadRequest("invalid day of week: %d", *params.DayOfWeek)
		}
		show.DayOfWeek = time.Weekday(*params.DayOfWeek)
	}

	if params.Website != nil {
		website := *params.Website
		if website != "" && !strings.HasPrefix(website, "http") {
			website = fmt.Sprintf("http://%s", website)
		}
		show.Website = website
	}

	if params.Wikipedia != nil {
		wikipedia := *params.Wikipedia
		if wikipedia != "" && !strings.HasPrefix(wikipedia, "http") {
			wikipedia = fmt.Sprintf("http://%s", wikipedia)
		}
		show.Wikipedia = wikipedia
	}

	if params.ExtraLink != nil {
		extraLink := *params.ExtraLink
		if extraLink != "" && !strings.HasPrefix(extraLink, "http") {
			extraLink = fmt.Sprintf("http://%s", extraLink)
		}
		show.ExtraLink = extraLink
	}

	if params.OnNetflix != nil {
		show.OnNetflix = *params.OnNetflix
	}

	if err = show.Save(); err != nil {
		return err
	}
	req.Ok(NewPayload(show))
	return nil
}
