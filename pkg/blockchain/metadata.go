//go:generate go run ../../resources/generate-smart-binds/main.go

package blockchain

import (
	"math/big"
	"matrix-ai-framework/pkg/db"
)

type OrganizationMetaData struct {
	SnetID      string  `json:"snet_id"` // id from blockchain
	ID          int     `json:"id"`      // internal id
	OrgName     string  `json:"org_name"`
	OrgID       string  `json:"org_id"`
	Groups      []Group `json:"groups"`
	OrgType     string  `json:"org_type"`
	Description struct {
		Description      string `json:"description"`
		ShortDescription string `json:"short_description"`
		URL              string `json:"url"`
	} `json:"description"`
	Contacts []struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		ContactType string `json:"contact_type"`
	} `json:"contacts"`
	Assets struct {
		HeroImage string `json:"hero_image"`
	} `json:"assets"`
	Owner string
}

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

type Group struct {
	GroupName        string   `json:"group_name"`
	GroupID          string   `json:"group_id"`
	PaymentDetails   Payment  `json:"payment"`
	LicenseEndpoints []string `json:"license_server_endpoints"`
}

type Payment struct {
	PaymentAddress              string                      `json:"payment_address"`
	PaymentExpirationThreshold  *big.Int                    `json:"payment_expiration_threshold"`
	PaymentChannelStorageType   string                      `json:"payment_channel_storage_type"`
	PaymentChannelStorageClient PaymentChannelStorageClient `json:"payment_channel_storage_client"`
}

type PaymentChannelStorageClient struct {
	ConnectionTimeout string   `json:"connection_timeout" mapstructure:"connection_timeout"`
	RequestTimeout    string   `json:"request_timeout" mapstructure:"request_timeout"`
	Endpoints         []string `json:"endpoints"`
}

type ServiceMetadata struct {
	ID            int
	SnetID        string
	SnetOrgID     string
	OrgID         int
	Version       int    `json:"version"`
	DisplayName   string `json:"display_name"`
	Encoding      string `json:"encoding"`
	ServiceType   string `json:"service_type"`
	ModelIpfsHash string `json:"model_ipfs_hash"`
	MpeAddress    string `json:"mpe_address"`
	Groups        []struct {
		FreeCalls             int      `json:"free_calls"`
		FreeCallSignerAddress string   `json:"free_call_signer_address"`
		DaemonAddresses       []string `json:"daemon_addresses"`
		Pricing               []struct {
			Default     bool   `json:"default"`
			PriceModel  string `json:"price_model"`
			PriceInCogs int    `json:"price_in_cogs"`
		} `json:"pricing"`
		Endpoints []string `json:"endpoints"`
		GroupID   string   `json:"group_id"`
		GroupName string   `json:"group_name"`
	} `json:"groups"`
	ServiceDescription struct {
		URL              string `json:"url"`
		ShortDescription string `json:"short_description"`
		Description      string `json:"description"`
	} `json:"service_description"`
	Media []struct {
		Order     int    `json:"order"`
		URL       string `json:"url"`
		FileType  string `json:"file_type"`
		AssetType string `json:"asset_type"`
		AltText   string `json:"alt_text"`
	} `json:"media"`
	Contributors []struct {
		Name    string `json:"name"`
		EmailID string `json:"email_id"`
	} `json:"contributors"`
	Tags []string `json:"tags"`
}
