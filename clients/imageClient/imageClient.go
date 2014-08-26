package imageClient

import (
	"fmt"
	"encoding/xml"
	"errors"
	azure "github.com/MSOpenTech/azure-sdk-for-go"
)

func GetImageList() (ImageList, error){
	imageList := ImageList{}

	requestURL :=  "services/images"
	response, err := azure.SendAzureGetRequest(requestURL)
	if err != nil {
		return imageList, err
	}

	err = xml.Unmarshal(response, &imageList)
	return imageList, err
}

func ResolveImageName(imageName string) (error) {
	imageList, err := GetImageList()
	if err != nil {
		return err
	}

	for _, image := range imageList.OSImages {
		if image.Name != imageName {
			continue
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Error: Can not find image %s in specified subscription, please specify another image name \n", imageName))
}
