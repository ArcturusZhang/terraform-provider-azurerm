package client

import (
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-12-01/armcompute"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/marketplaceordering/mgmt/2015-06-01/marketplaceordering"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AvailabilitySetsClient          *armcompute.AvailabilitySetsClient
	DedicatedHostsClient            *armcompute.DedicatedHostsClient
	DedicatedHostGroupsClient       *armcompute.DedicatedHostGroupsClient
	DisksClient                     *armcompute.DisksClient
	DiskAccessClient                *armcompute.DiskAccessesClient
	DiskEncryptionSetsClient        *compute.DiskEncryptionSetsClient
	GalleriesClient                 *compute.GalleriesClient
	GalleryImagesClient             *compute.GalleryImagesClient
	GalleryImageVersionsClient      *compute.GalleryImageVersionsClient
	ProximityPlacementGroupsClient  *armcompute.ProximityPlacementGroupsClient
	MarketplaceAgreementsClient     *marketplaceordering.MarketplaceAgreementsClient
	ImagesClient                    *compute.ImagesClient
	SnapshotsClient                 *armcompute.SnapshotsClient
	UsageClient                     *compute.UsageClient
	VMExtensionImageClient          *compute.VirtualMachineExtensionImagesClient
	VMExtensionClient               *compute.VirtualMachineExtensionsClient
	VMScaleSetClient                *compute.VirtualMachineScaleSetsClient
	VMScaleSetExtensionsClient      *compute.VirtualMachineScaleSetExtensionsClient
	VMScaleSetRollingUpgradesClient *compute.VirtualMachineScaleSetRollingUpgradesClient
	VMScaleSetVMsClient             *compute.VirtualMachineScaleSetVMsClient
	VMClient                        *compute.VirtualMachinesClient
	VMImageClient                   *compute.VirtualMachineImagesClient
	SSHPublicKeysClient             *compute.SSHPublicKeysClient
}

func NewClient(o *common.ClientOptions) *Client {
	diskEncryptionSetsClient := compute.NewDiskEncryptionSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&diskEncryptionSetsClient.Client, o.ResourceManagerAuthorizer)

	galleriesClient := compute.NewGalleriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleriesClient.Client, o.ResourceManagerAuthorizer)

	galleryImagesClient := compute.NewGalleryImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleryImagesClient.Client, o.ResourceManagerAuthorizer)

	galleryImageVersionsClient := compute.NewGalleryImageVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleryImageVersionsClient.Client, o.ResourceManagerAuthorizer)

	imagesClient := compute.NewImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&imagesClient.Client, o.ResourceManagerAuthorizer)

	marketplaceAgreementsClient := marketplaceordering.NewMarketplaceAgreementsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&marketplaceAgreementsClient.Client, o.ResourceManagerAuthorizer)

	usageClient := compute.NewUsageClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&usageClient.Client, o.ResourceManagerAuthorizer)

	vmExtensionImageClient := compute.NewVirtualMachineExtensionImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmExtensionImageClient.Client, o.ResourceManagerAuthorizer)

	vmExtensionClient := compute.NewVirtualMachineExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmExtensionClient.Client, o.ResourceManagerAuthorizer)

	vmImageClient := compute.NewVirtualMachineImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmImageClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetExtensionsClient := compute.NewVirtualMachineScaleSetExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetExtensionsClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetRollingUpgradesClient := compute.NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetRollingUpgradesClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetVMsClient := compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetVMsClient.Client, o.ResourceManagerAuthorizer)

	vmClient := compute.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmClient.Client, o.ResourceManagerAuthorizer)

	sshPublicKeysClient := compute.NewSSHPublicKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sshPublicKeysClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AvailabilitySetsClient:          armcompute.NewAvailabilitySetsClient(o.ResourceManagerConnection, o.SubscriptionId),
		DedicatedHostsClient:            armcompute.NewDedicatedHostsClient(o.ResourceManagerConnection, o.SubscriptionId),
		DedicatedHostGroupsClient:       armcompute.NewDedicatedHostGroupsClient(o.ResourceManagerConnection, o.SubscriptionId),
		DisksClient:                     armcompute.NewDisksClient(o.ResourceManagerConnection, o.SubscriptionId),
		DiskAccessClient:                armcompute.NewDiskAccessesClient(o.ResourceManagerConnection, o.SubscriptionId),
		DiskEncryptionSetsClient:        &diskEncryptionSetsClient,
		GalleriesClient:                 &galleriesClient,
		GalleryImagesClient:             &galleryImagesClient,
		GalleryImageVersionsClient:      &galleryImageVersionsClient,
		ImagesClient:                    &imagesClient,
		MarketplaceAgreementsClient:     &marketplaceAgreementsClient,
		ProximityPlacementGroupsClient:  armcompute.NewProximityPlacementGroupsClient(o.ResourceManagerConnection, o.SubscriptionId),
		SnapshotsClient:                 armcompute.NewSnapshotsClient(o.ResourceManagerConnection, o.SubscriptionId),
		UsageClient:                     &usageClient,
		VMExtensionImageClient:          &vmExtensionImageClient,
		VMExtensionClient:               &vmExtensionClient,
		VMScaleSetClient:                &vmScaleSetClient,
		VMScaleSetExtensionsClient:      &vmScaleSetExtensionsClient,
		VMScaleSetRollingUpgradesClient: &vmScaleSetRollingUpgradesClient,
		VMScaleSetVMsClient:             &vmScaleSetVMsClient,
		VMClient:                        &vmClient,
		VMImageClient:                   &vmImageClient,
		SSHPublicKeysClient:             &sshPublicKeysClient,
	}
}
