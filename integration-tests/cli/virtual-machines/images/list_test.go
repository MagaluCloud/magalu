package images

import "testing"

func TestListImages(t *testing.T) {
	imagesResponse, err := listImages()
	if err != nil {
		t.Errorf("Error listing images: %v", err)
	}

	if len(imagesResponse.Images) == 0 {
		t.Errorf("No images found")
	}

	for _, image := range imagesResponse.Images {
		if image.Name == "" {
			t.Errorf("Image name is empty")
		}
		if image.ID == "" {
			t.Errorf("Image ID is empty")
		}
	}
}
