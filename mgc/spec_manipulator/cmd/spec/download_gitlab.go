package spec

import (
	"fmt"
	"os"
	"strings"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

const projectID = 7739
const gitlabAPI = "https://gitlab.luizalabs.com/api/v4"

func validateGitlabUrl(url string) (bool, error) {
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	git, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabAPI))
	if err != nil {
		return false, fmt.Errorf("failed to create client: %v", err)
	}
	// file api_products/mcr-api/br-ne1-prod-yel-1/openapi.yaml
	sliceURL := strings.Split(url, "/")
	filePath := strings.Join(sliceURL[1:], "/")
	branch := sliceURL[0]

	_, _, err = git.RepositoryFiles.GetRawFile(projectID, filePath, &gitlab.GetRawFileOptions{
		Ref: &branch,
	})
	if err != nil {
		return false, fmt.Errorf("failed to get file: %v", err)
	}

	return true, nil
}

func downloadGitlab(url, caminhoDestino string) error {
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		return fmt.Errorf("GITLAB_TOKEN is not set")
	}

	git, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabAPI))
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// file api_products/mcr-api/br-ne1-prod-yel-1/openapi.yaml
	url = strings.Split(url, "?")[0]
	url = strings.TrimPrefix(url, "https://gitlab.luizalabs.com/open-platform/pcx/u0/-/blob/")

	sliceURL := strings.Split(url, "/")
	filePath := strings.Join(sliceURL[1:], "/")
	branch := sliceURL[0]

	// Obter o conteúdo do arquivo
	file, _, err := git.RepositoryFiles.GetRawFile(projectID, filePath, &gitlab.GetRawFileOptions{
		Ref: &branch,
	})
	if err != nil {
		return fmt.Errorf("failed to get file: %v", err)
	}

	if err := os.WriteFile(caminhoDestino, file, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}
