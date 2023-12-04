package main

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

func main() {
	fmt.Println("hi world, this is the gh-sarif extension!")

	// Get the repository's full name from the GITHUB_REPOSITORY environment variable
	// repoFullName := os.Getenv("GITHUB_REPOSITORY")
	// if repoFullName == "" {
	// 	fmt.Println("GITHUB_REPOSITORY environment variable is not set")
	// 	return
	// }

	// Parse the repository's full name to get a ghrepo.Interface
	// repo, err := ghrepo.FromFullName(repoFullName)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	repo, err := repository.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	response := struct {Login string}{}
	err = client.Get(fmt.Sprintf("repos/%s/%s/code-scanning/alerts", repo.Owner, repo.Name), &response)
	if err != nil {
		fmt.Println(err)
		fmt.Println("response: ", response)
		return
	}

	fmt.Printf("running as %s\n", response.Login)
}


// package main

// import (
// 	"fmt"

// 	"github.com/cli/go-gh/v2/pkg/api"
// )

// func main() {
// 	fmt.Println("hi world, this is the gh-sarif extension!")
// 	client, err := api.DefaultRESTClient()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	response := struct {Login string}{}
// 	err = client.Get("repos/OWNER/REPO/code-scanning/alerts", &response)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("running as %s\n", response.Login)
// }

// // For more examples of using go-gh, see:
// // https://github.com/cli/go-gh/blob/trunk/example_gh_test.go