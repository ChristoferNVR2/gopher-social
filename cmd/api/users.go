package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
	"github.com/ChristoferNVR2/gopher-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user by ID
//	@Description	Fetches a user profile by ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	models.UserWithRole
//	@Failure		400		{object}	error	"Invalid user ID"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error	"Failed to fetch the user"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Makes the given follower follow the user identified by the path ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int				true	"ID of the user to follow"
//	@Param			payload	body	FollowerPayload	true	"Follower payload"
//	@Success		204		"User followed"
//	@Failure		400		{object}	error	"Invalid user ID or malformed JSON"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		409		{object}	error	"Already following"
//	@Failure		500		{object}	error	"Failed to follow the user"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, followedID); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollows a user
//	@Description	Makes the given follower unfollow the user identified by the path ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int				true	"ID of the user to unfollow"
//	@Param			payload	body	FollowerPayload	true	"Follower payload"
//	@Success		204		"User unfollowed"
//	@Failure		400		{object}	error	"Invalid user ID or malformed JSON"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error	"Failed to unfollow the user"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, followedID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Registers a user account
//	@Description	Activates a user account using the provided activation token.
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Activation token"
//	@Success		204		{string}	string	"User account activated successfully"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	ctx := r.Context()

	err := app.store.Users.Activate(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *models.UserWithRole {
	user, _ := r.Context().Value(userCtx).(*models.UserWithRole)
	return user
}
