package osimage

import (
	"encoding/xml"

	"github.com/Azure/azure-sdk-for-go/management"
)

// OSImageClient is used to perform operations on Azure Locations
type OSImageClient struct {
	client management.Client
}

type ListOSImagesResponse struct {
	XMLName  xml.Name  `xml:"Images"`
	OSImages []OSImage `xml:"OSImage"`
}

type OSImage struct {
	XMLName           xml.Name `xml:"http://schemas.microsoft.com/windowsazure OSImage"`
	Category          string   `xml:",omitempty"` // Public || Private || MSDN
	Label             string   // Specifies an identifier for the image.
	MediaLink         string   // Specifies the location of the vhd file for the image. The storage account where the vhd is located must be associated with the specified subscription.
	Name              string   // Specifies the name of the operating system image. This is the name that is used when creating one or more virtual machines using the image.
	OS                string   // Linux || Windows
	Eula              string   `xml:",omitempty"` // Specifies the End User License Agreement that is associated with the image. The value for this element is a string, but it is recommended that the value be a URL that points to a EULA.
	Description       string   `xml:",omitempty"` // Specifies the description of the image.
	ImageFamily       string   `xml:",omitempty"` // Specifies a value that can be used to group images.
	PublishedDate     string   `xml:",omitempty"` // Specifies the date when the image was added to the image repository.
	ShowInGui         bool     // Specifies whether the image should appear in the image gallery.
	PrivacyURI        string   `xml:"PrivacyUri,omitempty"`   // Specifies the URI that points to a document that contains the privacy policy related to the image.
	IconURI           string   `xml:"IconUri,omitempty"`      // Specifies the Uri to the icon that is displayed for the image in the Management Portal.
	RecommendedVMSize string   `xml:",omitempty"`             // Specifies the size to use for the virtual machine that is created from the image.
	SmallIconURI      string   `xml:"SmallIconUri,omitempty"` // Specifies the URI to the small icon that is displayed when the image is presented in the Microsoft Azure Management Portal.
	Language          string   `xml:",omitempty"`             // Specifies the language of the image.

	LogicalSizeInGB   float64 `xml:",omitempty"` //Specifies the size, in GB, of the image.
	Location          string  `xml:",omitempty"` // The geo-location in which this media is located. The Location value is derived from storage account that contains the blob in which the media is located. If the storage account belongs to an affinity group the value is NULL.
	AffinityGroup     string  `xml:",omitempty"` // Specifies the affinity in which the media is located. The AffinityGroup value is derived from storage account that contains the blob in which the media is located. If the storage account does not belong to an affinity group the value is NULL and the element is not displayed in the response. This value is NULL for platform images.
	IsPremium         string  `xml:",omitempty"` // Indicates whether the image contains software or associated services that will incur charges above the core price for the virtual machine. For additional details, see the PricingDetailLink element.
	PublisherName     string  `xml:",omitempty"` // The name of the publisher of the image. All user images have a publisher name of User.
	PricingDetailLink string  `xml:",omitempty"` // Specifies a URL for an image with IsPremium set to true, which contains the pricing details for a virtual machine that is created from the image.
	IOType            IOType  `xml:",omitempty"` // Provisioned || Standard
}

type IOType string

const (
	IOTypeProvisioned IOType = "Provisioned"
	IOTypeStandard    IOType = "Standard"
)
