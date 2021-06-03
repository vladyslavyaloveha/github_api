package main

import (
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	_ "github_api/docs"
)

var CLIENT = getGithubClient(getEnvVar("GITHUB_TOKEN"))

// @title Fetching Github Example API
// @version 1.0
// @description This is a simple Github scrapper server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /github
func main() {
	app := iris.New()
	docsAPI := app.Party("/swagger")
	{
		docsAPI.Use(iris.Compression)
		swaggerUI := swagger.WrapHandler(swaggerFiles.Handler)
		// Register on http://localhost:8080/swagger
		docsAPI.Get("/swagger", swaggerUI)
		// And the wildcard one for index.html, *.js, *.css and e.t.c.
		app.Get("/swagger/{any:path}", swaggerUI)
	}

	githubAPI := app.Party("/github")
	{
		githubAPI.Use(iris.Compression)
		// GET: http://localhost:8080/github/repositories
		githubAPI.Get("/repositories", listOwnersRepositories)
		// GET: http://localhost:8080/github/{owner}/{name}
		githubAPI.Get("/{owner:string}/{name:string}", listRepository)
		// GET: http://localhost:8080/github/{name}/issues
		githubAPI.Get("/{name:string}/issues/", listIssues)
		// GET: http://localhost:8080/github/{name}/commits
		githubAPI.Get("/{name:string}/commits", listCommits)
	}
	app.Listen(":8080")
}
