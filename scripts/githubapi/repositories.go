package githubapi

import "fmt"

const reposPerPage = 100

// FetchRepositories fetches all owned repositories for a user using cursor-based pagination.
// maxRepos limits the total number of repositories fetched (0 = no limit beyond API constraints).
func FetchRepositories(token, username string, maxRepos int) ([]RepoNode, error) {
	if maxRepos <= 0 {
		maxRepos = reposPerPage
	}
	perPage := reposPerPage
	if maxRepos < perPage {
		perPage = maxRepos
	}

	firstQuery := fmt.Sprintf(`
		query ($login: String!) {
			user(login: $login) {
				repositories(first: %d, ownerAffiliations: OWNER) {
					edges { cursor }
					nodes { forkCount stargazerCount }
				}
			}
		}
	`, perPage)

	result, err := GqlQuery(token, firstQuery, map[string]any{"login": username})
	if err != nil {
		return nil, err
	}
	if errs := gqlErrors(result); errs != "" {
		return nil, fmt.Errorf("GraphQL error: %s", errs)
	}

	nodes, edges := extractRepoPage(result)
	allNodes := append([]RepoNode{}, nodes...)

	for len(nodes) == perPage && len(edges) > 0 && len(allNodes) < maxRepos {
		cursor := edges[len(edges)-1]
		nextNodes, nextEdges, err := fetchNextRepoPage(token, username, cursor, perPage)
		if err != nil {
			return nil, err
		}
		if len(nextNodes) == 0 {
			break
		}
		allNodes = append(allNodes, nextNodes...)
		nodes = nextNodes
		edges = nextEdges
	}

	return allNodes, nil
}

func fetchNextRepoPage(token, username, cursor string, perPage int) ([]RepoNode, []string, error) {
	query := fmt.Sprintf(`
		query ($login: String!, $cursor: String!) {
			user(login: $login) {
				repositories(after: $cursor, first: %d, ownerAffiliations: OWNER) {
					edges { cursor }
					nodes { forkCount stargazerCount }
				}
			}
		}
	`, perPage)

	result, err := GqlQuery(token, query, map[string]any{"login": username, "cursor": cursor})
	if err != nil {
		return nil, nil, err
	}
	if errs := gqlErrors(result); errs != "" {
		return nil, nil, fmt.Errorf("GraphQL error (next page): %s", errs)
	}

	nodes, edges := extractRepoPage(result)
	return nodes, edges, nil
}

// extractRepoPage pulls nodes and edge cursors from a GraphQL repositories response.
func extractRepoPage(result map[string]any) ([]RepoNode, []string) {
	data, _ := result["data"].(map[string]any)
	user, _ := data["user"].(map[string]any)
	repos, _ := user["repositories"].(map[string]any)

	rawNodes, _ := repos["nodes"].([]any)
	nodes := make([]RepoNode, 0, len(rawNodes))
	for _, n := range rawNodes {
		node, _ := n.(map[string]any)
		nodes = append(nodes, RepoNode{
			ForkCount:      AsInt(node, "forkCount"),
			StargazerCount: AsInt(node, "stargazerCount"),
		})
	}

	rawEdges, _ := repos["edges"].([]any)
	edges := make([]string, 0, len(rawEdges))
	for _, e := range rawEdges {
		edge, _ := e.(map[string]any)
		cursor, _ := edge["cursor"].(string)
		edges = append(edges, cursor)
	}

	return nodes, edges
}

// gqlErrors returns the first error message from a GraphQL response, or empty string.
func gqlErrors(result map[string]any) string {
	errs, _ := result["errors"].([]any)
	if len(errs) == 0 {
		return ""
	}
	first, _ := errs[0].(map[string]any)
	msg, _ := first["message"].(string)
	return msg
}
