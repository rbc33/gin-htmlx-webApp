/*
	    There's a lot of common functionality between the image
		and gallery handlers, so this file is meant to share those
		functionalities.
*/

package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// This list contains the valid file
// extensions for an image.
var ValidImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

func populateImageMetadata(metadata_path string) (Image, error) {

	// Check if a json metadata file exists
	metadata_contents, err := os.ReadFile(path.Join(Settings.ImageDirectory, metadata_path))
	if err != nil {
		return Image{}, fmt.Errorf("could not read metadata for image `%s`", metadata_path)
	}

	var image Image
	err = json.Unmarshal(metadata_contents, &image)

	if err != nil {
		return Image{}, fmt.Errorf("could not deserailize metadata for image `%s`", metadata_path)
	}

	ext := path.Ext(image.Filename)

	// Checking for the existence of a value in a map takes O(1) and therefore it's faster than
	// iterating over a string slice
	_, ok := ValidImageExtensions[ext]
	if !ok {
		return Image{}, fmt.Errorf("image type provided in metadata `%s` is not supported: `%s`", metadata_path, image.Filename)
	}

	filepath := path.Join("/images/data/", image.Filename)

	metadata_uuid := strings.TrimSuffix(metadata_path, ext)
	image.Ext = ext
	image.Uuid = metadata_uuid
	image.Filepath = filepath
	return image, nil
}

// Given a list of files, this function will return
// a filtered list of valid images, with the page number
// and page size taken as pagination arguments.
//
// paths must be a list of strings referencing the metadata file for an image.
// page_size must be a non-negative number greater than zero.
// page_num must be a non-negative number greater than 0.
func GetImages(files []os.DirEntry, page_size, page_num int) ([]Image, error) {

	limit := page_size
	offset := (page_num - 1) * page_size

	// Get all the files inside the image directory
	limit = min(offset+limit, len(files))
	offset = min(offset, len(files))

	// Filter all the non-images out of the list
	valid_images := make([]Image, 0)
	for _, file := range files[offset:limit] {

		filename := file.Name()
		ext := path.Ext(file.Name())
		// Checking for the existence of a value in a map takes O(1) and therefore it's faster than
		// iterating over a string slice
		_, ok := ValidImageExtensions[ext]
		if !ok {
			continue
		}

		image := Image{
			Uuid: filename[:len(filename)-len(ext)],
			Name: filename,
			Ext:  ext,
		}
		valid_images = append(valid_images, image)
	}

	return valid_images, nil
}
