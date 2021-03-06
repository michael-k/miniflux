// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui  // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/context"
	"miniflux.app/http/request"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/model"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
)

// ShowSearchEntries shows all entries for the given feed.
func (c *Controller) ShowSearchEntries(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	searchQuery := request.QueryParam(r, "q", "")
	offset := request.QueryIntParam(r, "offset", 0)
	builder := c.store.NewEntryQueryBuilder(user.ID)
	builder.WithSearchQuery(searchQuery)
	builder.WithoutStatus(model.EntryStatusRemoved)
	builder.WithOrder(model.DefaultSortingOrder)
	builder.WithDirection(user.EntryDirection)
	builder.WithOffset(offset)
	builder.WithLimit(nbItemsPerPage)

	entries, err := builder.GetEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	count, err := builder.CountEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	pagination := c.getPagination(route.Path(c.router, "searchEntries"), count, offset)
	pagination.SearchQuery = searchQuery

	view.Set("searchQuery", searchQuery)
	view.Set("entries", entries)
	view.Set("total", count)
	view.Set("pagination", pagination)
	view.Set("menu", "search")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))
	view.Set("hasSaveEntry", c.store.HasSaveEntry(user.ID))

	html.OK(w, r, view.Render("search_entries"))
}
