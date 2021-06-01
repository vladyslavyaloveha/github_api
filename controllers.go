package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
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
	fmt.Println(repository)
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
// @Param until query string false "Since timestamp: 2020-05-25T06:34:16Z"
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
