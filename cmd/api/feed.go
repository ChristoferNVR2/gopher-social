package main

import (
	"net/http"

	"github.com/ChristoferNVR2/gopher-social/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the posts of the user and the users they follow, with pagination, sorting and filtering.
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Page size (1-20)"		default(20)
//	@Param			offset	query		int		false	"Offset"				default(0)
//	@Param			sort	query		string	false	"Sort by created_at"	Enums(asc, desc)	default(desc)
//	@Param			tags	query		string	false	"Comma-separated tags (max 5)"
//	@Param			search	query		string	false	"Search in title and content (max 100 chars)"
//	@Param			since	query		string	false	"Created since (RFC3339)"
//	@Param			until	query		string	false	"Created until (RFC3339)"
//	@Success		200		{array}		[]models.PostWithMetadata
//	@Failure		400		{object}	error	"Invalid query params (bad limit, offset, sort, tags, search or date)"
//	@Failure		500		{object}	error	"Failed to fetch the feed"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: pagination, filtering, etc.
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(66), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
