package autoruns

import (
	"github.com/botherder/go-autoruns"
)

type AutorunsInfo struct {
	Type      string `json:"type"`
	Location  string `json:"location"`
	ImagePath string `json:"image_path"`
	ImageName string `json:"image_name"`
	Arguments string `json:"arguments"`
	MD5       string `json:"md5"`
	SHA1      string `json:"sha1"`
	SHA256    string `json:"sha256"`
}

func GetAutoruns() []AutorunsInfo {

	autoruns := autoruns.Autoruns()
	autoslice := make([]AutorunsInfo, 0)

	for _, autorun := range autoruns {

		helium := AutorunsInfo{
			Type:      autorun.Type,
			Location:  autorun.Location,
			ImagePath: autorun.ImagePath,
			ImageName: autorun.ImageName,
			Arguments: autorun.Arguments,
			MD5:       autorun.MD5,
			SHA1:      autorun.SHA1,
			SHA256:    autorun.SHA256}
		autoslice = append(autoslice, helium)

	}

	// define a map to keep track of seen image paths
	seenPaths := make(map[string]bool)

	// define a new array to store unique items
	uniqueItems := []AutorunsInfo{}

	// iterate over the original array
	for _, item := range autoslice {
		// check if the image path has already been seen
		if !seenPaths[item.ImagePath] {
			// add the item to the new array
			uniqueItems = append(uniqueItems, item)
		}
		// mark the image path as seen
		seenPaths[item.ImagePath] = true
	}

	return (uniqueItems)
}
