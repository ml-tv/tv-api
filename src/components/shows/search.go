package shows

import (
	"fmt"

	"strings"

	"strconv"

	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/paginator"
	"github.com/ml-tv/tv-api/src/core/primitives/slices"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/storage/db"
	"github.com/ml-tv/tv-api/src/services/tmdb"
)

// DefaultNbResultsPerPage represents the default number of result per page
const DefaultNbResultsPerPage = 20

// SearchParams represents the params needed by the Search handler
type SearchParams struct {
	paginator.HandlerParams

	// FromProvider will query TMDb instead of the local database
	FromProvider bool `from:"query" json:"from_provider" default:"false"`

	// Name represents a string to use to look against the name field
	Name string `from:"query" json:"name" params:"trim"`

	// Status represents a list of status separated by "|"
	// ex ?status=0|1|3
	Status string `from:"query" json:"status" params:"trim"`

	// DayOfWeek filters the day of the week
	DayOfWeek *int `from:"query" json:"day_of_week"`

	// OrderBy represents a list of orders separated by "|"
	// ex ?order=-name|day_of_week will order by name desc and day of week asc
	// fields are: name, day_of_week, returning_date, created_at
	OrderBy string `from:"query" json:"order" params:"trim"`
}

// searchTMDb is an API handler to search a show on TMDb
func searchTMDb(req *router.Request, params *SearchParams) error {
	pagination := params.Paginator(DefaultNbResultsPerPage)
	if !pagination.IsValid() {
		return httperr.NewBadRequest("invalid pagination data")
	}
	if params.Name == "" {
		return httperr.NewBadRequest("you need to provide a name")
	}

	shows, err := tmdb.SearchShows(params.Name, pagination.CurrentPage)
	if err != nil {
		return err
	}

	s, err := NewListFromTMDb(shows)
	if err != nil {
		return err
	}

	req.Ok(NewPayloadList(s))
	return nil
}

// Search is an API handler to search a show
func Search(req *router.Request) error {
	params := req.Params.(*SearchParams)

	if params.FromProvider {
		return searchTMDb(req, params)
	}

	pagination := params.Paginator(DefaultNbResultsPerPage)
	if !pagination.IsValid() {
		return httperr.NewBadRequest("Invalid pagination data")
	}

	// Set default SQL params
	selct := "*"
	from := "shows"
	where := ""
	orderBy := ""
	args := map[string]interface{}{}

	// Full text search on the name
	if params.Name != "" {
		selct += ", ts_rank(name_vector, keywords, 1) AS rank"
		from += fmt.Sprintf(", plainto_tsquery(:name) keywords")
		where += "keywords @@ name_vector"
		orderBy = "rank DESC"
		args["name"] = params.Name
	}

	// Filter the Day of week
	if params.DayOfWeek != nil {
		if *params.DayOfWeek < 0 || *params.DayOfWeek > 6 {
			return httperr.NewBadRequest("day of week must be between 0 and 6")
		}
		if where != "" {
			where += " AND "
		}
		where += "day_of_week=:dow"
		args["dow"] = *params.DayOfWeek
	}

	// Filter the Status
	if params.Status != "" {
		fields := strings.Split(params.Status, "|")
		// "in" contains what's should go inside the parenthesis in "IN()"
		// "in" always starts with a comma, so it needs to be trimmed
		// ex: ,:status0,:status3
		in := ""
		for i, f := range fields {
			n, err := strconv.Atoi(f)
			if err != nil || n < 0 || n >= ShowStatusEndOfList {
				return httperr.NewBadRequest("not a valid status: %s", f)
			}
			// we cannot do something like IN(:statuses), so instead we do
			// IN(:status0,:status1,:status2,:status3,:statusX)
			sqlName := fmt.Sprintf("status%d", i)
			in += fmt.Sprintf(",:%s", sqlName)
			args[sqlName] = n
		}
		if in != "" {
			if where != "" {
				where += " AND "
			}
			where += fmt.Sprintf("status IN(%s)", in[1:])
		}
	}

	// Set sorting
	if params.OrderBy != "" {
		sortableFields := []string{"name", "day_of_week", "returning_date", "created_at"}
		fields := strings.Split(params.OrderBy, "|")
		for _, f := range fields {
			f = strings.ToLower(f)
			// we need at least 2 chars (ex. -a)
			if len(f) < 2 {
				return httperr.NewBadRequest("invalid sort option: %s", f)
			}
			order := "ASC"
			if f[0] == '-' {
				order = "DESC"
				f = f[1:]
			}
			found, err := slices.InSlice(sortableFields, f)
			if err != nil {
				return err
			}
			if !found {
				return httperr.NewBadRequest("field not sortable: %s", f)
			}
			if orderBy != "" {
				orderBy += ", "
			}
			orderBy += fmt.Sprintf("%s %s", f, order)
		}
	}

	// Set SQL keywords if needed
	if where != "" {
		where = " WHERE " + where
	}
	// We default the ordering by recency
	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	// Exec query and return payload
	var list []*Show
	stmt := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY %s LIMIT %d OFFSET %d",
		selct, from, where, orderBy, pagination.Limit(), pagination.Offset())

	// stmt := fmt.Sprintf("SELECT %s FROM %s %s %s", selct, from, where, orderBy)
	if err := db.NamedSelect(&list, stmt, args); err != nil {
		return err
	}
	req.Ok(NewPayloadList(list))
	return nil
}
