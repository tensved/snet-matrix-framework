//go:generate go run ../../resources/generate-smart-binds/main.go

package blockchain

import (
	"github.com/tensved/snet-matrix-framework/pkg/db"
	"math/big"
)

// OrganizationMetaData represents metadata for an organization in the blockchain.
type OrganizationMetaData struct {
	SnetID      string  `json:"snet_id"`  // ID from blockchain.
	ID          int     `json:"id"`       // Internal ID.
	OrgName     string  `json:"org_name"` // Organization name.
	OrgID       string  `json:"org_id"`   // Organization ID.
	Groups      []Group `json:"groups"`   // Groups within the organization.
	OrgType     string  `json:"org_type"` // Organization type.
	Description struct {
		Description      string `json:"description"`       // Detailed description.
		ShortDescription string `json:"short_description"` // Short description.
		URL              string `json:"url"`               // URL.
	} `json:"description"` // Organization description.
	Contacts []struct {
		Email       string `json:"email"`        // Contact email.
		Phone       string `json:"phone"`        // Contact phone.
		ContactType string `json:"contact_type"` // Type of contact.
	} `json:"contacts"` // List of contacts.
	Assets struct {
		HeroImage string `json:"hero_image"` // Hero image URL.
	} `json:"assets"` // Organization assets.
	Owner string // Owner of the organization.
}

// DB converts OrganizationMetaData to db.SnetOrganization and a slice of db.SnetOrgGroup.
//
// Returns:
//   - db.SnetOrganization: The database representation of the organization.
//   - []db.SnetOrgGroup: A slice of organization groups in the database format.
func (o OrganizationMetaData) DB() (db.SnetOrganization, []db.SnetOrgGroup) {
	org := db.SnetOrganization{
		SnetID:           o.SnetID,
		ID:               o.ID,
		Name:             o.OrgName,
		Type:             o.OrgType,
		Description:      o.Description.Description,
		ShortDescription: o.Description.ShortDescription,
		URL:              o.Description.URL,
		Owner:            o.Owner,
		Image:            o.Assets.HeroImage,
	}
	var groups []db.SnetOrgGroup
	for _, group := range o.Groups {
		groups = append(groups, db.SnetOrgGroup{
			GroupID:                    group.GroupID,
			GroupName:                  group.GroupName,
			PaymentAddress:             group.PaymentDetails.PaymentAddress,
			PaymentExpirationThreshold: group.PaymentDetails.PaymentExpirationThreshold,
		})
	}
	return org, groups
}

// DB converts ServiceMetadata to db.SnetService.
//
// Returns:
//   - db.SnetService: The database representation of the service.
func (s ServiceMetadata) DB() db.SnetService {
	return db.SnetService{
		ID:                    s.ID,
		SnetID:                s.SnetID,
		SnetOrgID:             s.SnetOrgID,
		OrgID:                 s.OrgID,
		Version:               s.Version,
		DisplayName:           s.DisplayName,
		Encoding:              s.Encoding,
		ServiceType:           s.ServiceType,
		ModelIpfsHash:         s.ModelIpfsHash,
		MPEAddress:            s.MpeAddress,
		URL:                   s.Groups[0].Endpoints[0],
		Price:                 s.Groups[0].Pricing[0].PriceInCogs,
		GroupID:               s.Groups[0].GroupID,
		FreeCalls:             s.Groups[0].FreeCalls,
		FreeCallSignerAddress: s.Groups[0].FreeCallSignerAddress,
		Description:           s.ServiceDescription.Description,
		ShortDescription:      s.ServiceDescription.ShortDescription,
	}
}

// Group represents a group within an organization in the blockchain.
type Group struct {
	GroupName        string   `json:"group_name"`               // Name of the group.
	GroupID          string   `json:"group_id"`                 // ID of the group.
	PaymentDetails   Payment  `json:"payment"`                  // Payment details for the group.
	LicenseEndpoints []string `json:"license_server_endpoints"` // License server endpoints.
}

// Payment represents the payment details for a group.
type Payment struct {
	PaymentAddress              string                      `json:"payment_address"`                // Payment address.
	PaymentExpirationThreshold  *big.Int                    `json:"payment_expiration_threshold"`   // Payment expiration threshold.
	PaymentChannelStorageType   string                      `json:"payment_channel_storage_type"`   // Payment channel storage type.
	PaymentChannelStorageClient PaymentChannelStorageClient `json:"payment_channel_storage_client"` // Payment channel storage client.
}

// PaymentChannelStorageClient represents the client for storing payment channel information.
type PaymentChannelStorageClient struct {
	ConnectionTimeout string   `json:"connection_timeout" mapstructure:"connection_timeout"` // Connection timeout.
	RequestTimeout    string   `json:"request_timeout" mapstructure:"request_timeout"`       // Request timeout.
	Endpoints         []string `json:"endpoints"`                                            // List of endpoints.
}

// ServiceMetadata represents metadata for a service in the blockchain.
type ServiceMetadata struct {
	ID            int    // Internal ID.
	SnetID        string // ID from blockchain.
	SnetOrgID     string // Organization ID from blockchain.
	OrgID         int    // Internal organization ID.
	Version       int    `json:"version"`         // Version of the service.
	DisplayName   string `json:"display_name"`    // Display name of the service.
	Encoding      string `json:"encoding"`        // Encoding type of the service.
	ServiceType   string `json:"service_type"`    // Type of the service.
	ModelIpfsHash string `json:"model_ipfs_hash"` // IPFS hash of the service model.
	MpeAddress    string `json:"mpe_address"`     // MPE address of the service.
	Groups        []struct {
		FreeCalls             int      `json:"free_calls"`               // Number of free calls.
		FreeCallSignerAddress string   `json:"free_call_signer_address"` // Address of the free call signer.
		DaemonAddresses       []string `json:"daemon_addresses"`         // Daemon addresses.
		Pricing               []struct {
			Default     bool   `json:"default"`       // Indicates if this is the default pricing.
			PriceModel  string `json:"price_model"`   // Pricing model.
			PriceInCogs int    `json:"price_in_cogs"` // Price in cogs.
		} `json:"pricing"` // Pricing details.
		Endpoints []string `json:"endpoints"`  // Endpoints for the group.
		GroupID   string   `json:"group_id"`   // ID of the group.
		GroupName string   `json:"group_name"` // Name of the group.
	} `json:"groups"` // List of groups.
	ServiceDescription struct {
		URL              string `json:"url"`               // URL of the service.
		ShortDescription string `json:"short_description"` // Short description of the service.
		Description      string `json:"description"`       // Detailed description of the service.
	} `json:"service_description"` // Service description.
	Media []struct {
		Order     int    `json:"order"`      // Order of the media.
		URL       string `json:"url"`        // URL of the media.
		FileType  string `json:"file_type"`  // File type of the media.
		AssetType string `json:"asset_type"` // Asset type of the media.
		AltText   string `json:"alt_text"`   // Alternative text for the media.
	} `json:"media"` // List of media items.
	Contributors []struct {
		Name    string `json:"name"`     // Name of the contributor.
		EmailID string `json:"email_id"` // Email ID of the contributor.
	} `json:"contributors"` // List of contributors.
	Tags []string `json:"tags"` // List of tags.
}
