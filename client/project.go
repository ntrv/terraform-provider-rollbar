/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package client

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// ListProjects queries API for the list of projects
func (c *RollbarApiClient) ListProjects() ([]Project, error) {
	path := "/api/1/projects"
	l := log.With().
		Str("path", path).
		Logger()

	url := c.url
	url.Path = path

	resp, err := c.resty.R().
		SetResult(ProjectListResult{}).
		SetError(ErrorResult{}).
		Get(url.String())
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		errResp := resp.Error().(*ErrorResult)
		l.Err(errResp).Send()
		return nil, errResp
	}
	lpr := resp.Result().(*ProjectListResult)

	// FIXME: After deleting a project through the API, it still shows up in
	//  the list of projects returned by the API - only with its name set to
	//  nil. This seemingly undesirable behavior should be fixed on the API
	//  side. We work around it by removing any result with an empty name.
	cleaned := make([]Project, 0)
	for _, proj := range lpr.Result {
		if proj.Name != "" {
			cleaned = append(cleaned, proj)
		}
	}

	return cleaned, nil
}

// CreateProject creates a new project
func (c *RollbarApiClient) CreateProject(name string) (*Project, error) {
	p := "/api/1/projects"
	l := log.With().
		Str("name", name).
		Str("path", p).
		Logger()
	l.Debug().Msg("Creating new project")

	u := *c.url
	u.Path = p

	resp, err := c.resty.R().
		SetBody(map[string]interface{}{"name": name}).
		SetResult(ProjectResult{}).
		SetError(ErrorResult{}).
		Post(u.String())
	l.Debug().Bytes("body", resp.Body()).Msg("Response Body")
	if err != nil {
		l.Err(err).Msg("Error creating project")
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return nil, er
	}

	pr := resp.Result().(*ProjectResult)
	if pr.Err != 0 {
		l.Error().Msg("Unexpected error creating project")
	}

	return &pr.Result, nil
}

// ReadProject fetches data for the specified Project from the Rollbar API.
func (c *RollbarApiClient) ReadProject(id int) (*Project, error) {
	p := fmt.Sprintf("/api/1/project/%v", id)

	l := log.With().
		Int("id", id).
		Str("path", p).
		Logger()
	l.Debug().Msg("Reading project from API")

	u := *c.url
	u.Path = p

	//c.resty.SetDebug(true)
	//rzl := RestyZeroLogger{l}
	//c.resty.SetLogger(rzl)

	resp, err := c.resty.R().
		SetResult(ProjectResult{}).
		SetError(ErrorResult{}).
		Get(u.String())
	if err != nil {
		l.Err(err).Msg("Error reading project")
		return nil, err
	}
	l.Debug().Bytes("body", resp.Body()).Msg("Response Body")
	if resp.StatusCode() != http.StatusOK {
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return nil, er
	}

	pr := resp.Result().(*ProjectResult)
	l.Debug().Interface("ProjectResult", pr).Send()
	if pr.Err != 0 {
		l.Error().Msg("Unexpected error reading project")
	}

	return &pr.Result, nil
}
