package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"strconv"
	"time"
)

var JSONOptions = iris.JSON{Secure: true}

// listRepository godoc
// @Summary Retrieves repository info based on given name
// @Produce json
// @Param owner path string true "Repository Owner"
// @Param name path string true "Repository Name"
// @Success 200 {string} ok
// @Router /{owner}/{name} [get]
func listRepository(ctx iris.Context) {
	repoName := ctx.Params().Get("name")
	owner := ctx.Params().Get("owner")
	repository, err := getRepository(CLIENT, owner+"/"+repoName)
	if err != nil {
		ctx.JSON(context.Map{"repository": repository, "error": err.Error()}, JSONOptions)
	} else {
		ctx.JSON(context.Map{"repository": repository, "error": nil}, JSONOptions)
	}
}

// listIssues godoc
// @Summary Retrieves issues on given repository name
// @Produce json
// @Param owner path string true "Repository Owner"
// @Param name path string true "Repository Name"
// @Param state query string false "Issue State" Enums(all, open, closed)
// @Success 200 {string} ok
// @Router /{owner}/{name}/issues [get]
func listIssues(ctx iris.Context) {
	repoName := ctx.Params().Get("name")
	owner := ctx.Params().Get("owner")
	state, err := getIssueState(ctx)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	issues, err := getIssues(CLIENT, owner+"/"+repoName, IssueFilter{State: state})
	if err != nil {
		ctx.JSON(context.Map{"issues": issues, "error": err.Error()}, JSONOptions)
	} else {
		ctx.JSON(context.Map{"issues": issues, "error": nil}, JSONOptions)
	}
}

// listCommits godoc
// @Summary Retrieves commits based on given repository name
// @Produce json
// @Param owner path string true "Repository Owner"
// @Param name path string true "Repository Name"
// @Param author query string false "Commit Author: GitHub login or email address"
// @Param since query string false "Since timestamp: 2020-05-25T06:34:16Z"
// @Param until query string false "Until timestamp: 2020-05-25T06:34:16Z"
// @Success 200 {string} ok
// @Router /{owner}/{name}/commits [get]
func listCommits(ctx iris.Context) {
	repoName := ctx.Params().Get("name")
	owner := ctx.Params().Get("owner")
	sinceParam := ctx.URLParamDefault("since", "")
	since, err := time.Parse(time.RFC3339, sinceParam)
	if sinceParam != "" && err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	untilParam := ctx.URLParamDefault("until", "")
	until, err := time.Parse(time.RFC3339, untilParam)
	if untilParam != "" && err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	commits, err := getCommits(CLIENT, owner+"/"+repoName, CommitFilter{Author: getCommitAuthor(ctx),
		Since: since, Until: until})
	if err != nil {
		ctx.JSON(context.Map{"commits": commits, "error": err.Error()}, JSONOptions)
	} else {
		ctx.JSON(context.Map{"commits": commits, "error": nil}, JSONOptions)
	}
}

// listOwnersRepositories godoc
// @Summary Retrieves repositories based on given owners names
// @Produce json
// @Param owners query string true "Pass json with keyword owners and array with values"
// @Param max_requests query integer false "Max concurrent requests (1-100), default=5" minimum(1) maximum(100)
// @Success 200 {string} ok
// @Router /repositories [get]
func listOwnersRepositories(ctx iris.Context) {
	ownersParam := ctx.URLParamDefault("owners", "")
	maxRequestsParam := ctx.URLParamDefault("max_requests", "5")
	maxRequests, _ := strconv.Atoi(maxRequestsParam)

	var val map[string][]string
	if err := json.Unmarshal([]byte(ownersParam), &val); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	if _, ok := val["owners"]; !ok {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("key owners not passed"))
		return
	}
	repositories := getOwnersRepositories(val["owners"], maxRequests)
	ctx.JSON(context.Map{"repositories": repositories, "error": nil}, JSONOptions)

}
