package images

import (
	"encoding/json"
	"os/exec"

	model "integration-tests/cli/virtual-machines"
)

func listImages() (model.ImageResponse, error) {
	cmd := exec.Command("mgc", "vm", "images", "list", "-o", "json", "--raw")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return model.ImageResponse{}, err
	}

	var imagesResponse model.ImageResponse

	err = json.Unmarshal(output, &imagesResponse)
	if err != nil {
		return model.ImageResponse{}, err
	}
	return imagesResponse, nil
}
