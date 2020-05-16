package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type Pagination struct {
	page    int
	perPage int
	offset  int
	limit   int
	start   int
	end     int
}

func getPaginationFromRequest(r *http.Request) (*Pagination, error) {
	pagination := Pagination{}
	pagination.perPage = 10
	var err error
	pageString := r.FormValue("page")
	if pageString != "" {
		pagination.page, err = strconv.Atoi(pageString)
		if err != nil {
			return &Pagination{}, err
		}
	} else {
		pagination.page = 1
	}
	perPageString := r.FormValue("perPage")
	if perPageString != "" {
		pagination.perPage, err = strconv.Atoi(perPageString)
		if err != nil {
			return &Pagination{}, err
		}
	}
	offsetString := r.FormValue("offset")
	if offsetString != "" {
		pagination.offset, err = strconv.Atoi(offsetString)
		if err != nil {
			return &Pagination{}, err
		}
	}
	limitString := r.FormValue("limit")
	if limitString != "" {
		pagination.limit, err = strconv.Atoi(limitString)
		if err != nil {
			return &Pagination{}, err
		}
	}
	start := pagination.perPage * (pagination.page - 1)
	end := pagination.perPage * (pagination.page)
	if pagination.limit != 0 && end > pagination.limit {
		end = pagination.limit
	}
	if pagination.perPage > pagination.limit && pagination.limit > 0 {
		pagination.perPage = pagination.limit
	}
	pagination.start = start
	pagination.end = end
	return &pagination, nil
}

func setPaginationHeaders(w http.ResponseWriter, r *http.Request, p *Pagination, total int) {
	w.Header().Add("limit", fmt.Sprintf("%d", p.limit))
	w.Header().Add("page", fmt.Sprintf("%d", p.page))
	w.Header().Add("perPage", fmt.Sprintf("%d", p.perPage))
	w.Header().Add("offset", fmt.Sprintf("%d", p.offset))

	links := calculatePaginationHeaders(r, p, total)

	linkHeaderValue := strings.Join(links, ",")
	w.Header().Set("Link", linkHeaderValue)
}

func calculatePaginationHeaders(r *http.Request, p *Pagination, total int) []string {
	var links []string
	parsedUrl, _ := url.Parse(r.URL.String())
	var scheme string
	if r.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}
	q := parsedUrl.Query()
	pageString := q["page"]
	currentPage := 1
	if len(pageString) > 0 {
		currentPage, _ = strconv.Atoi(pageString[0])
	}
	nextPage := currentPage + 1
	previousPage := currentPage - 1
	if total < p.limit || p.limit == 0 {
		p.limit = total
	}
	lastPage := p.limit / p.perPage
	if p.limit%p.perPage > 0 {
		lastPage++
	}

	q.Set("page", fmt.Sprintf("%d", currentPage))
	parsedUrl.RawQuery = q.Encode()
	current := scheme + "://" + path.Join(r.Host, parsedUrl.String())

	q.Set("page", fmt.Sprintf("%d", 1))
	parsedUrl.RawQuery = q.Encode()
	first := scheme + "://" + path.Join(r.Host, parsedUrl.String())

	q.Set("page", fmt.Sprintf("%d", nextPage))
	parsedUrl.RawQuery = q.Encode()
	next := scheme + "://" + path.Join(r.Host, parsedUrl.String())

	q.Set("page", fmt.Sprintf("%d", previousPage))
	parsedUrl.RawQuery = q.Encode()
	previous := scheme + "://" + path.Join(r.Host, parsedUrl.String())

	q.Set("page", fmt.Sprintf("%d", lastPage))
	parsedUrl.RawQuery = q.Encode()
	last := scheme + "://" + path.Join(r.Host, parsedUrl.String())

	links = append(links, fmt.Sprintf(`<%s>; rel="%s"`, current, "self"))
	links = append(links, fmt.Sprintf(`<%s>; rel="%s"`, first, "first"))
	if nextPage <= lastPage {
		links = append(links, fmt.Sprintf(`<%s>; rel="%s"`, next, "next"))
	}
	if previousPage != 0 {
		links = append(links, fmt.Sprintf(`<%s>; rel="%s"`, previous, "prev"))
	}
	links = append(links, fmt.Sprintf(`<%s>; rel="%s"`, last, "last"))

	return links
}
