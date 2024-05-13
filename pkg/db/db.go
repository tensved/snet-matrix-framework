package db

import (
	"math/big"
	"time"
)

type Service interface {
	GetSnetOrgs() ([]SnetOrganization, error)
	GetSnetServices() ([]SnetService, error)
	GetSnetService(snetID string) (s SnetService, err error)
	CreateSnetService(service SnetService) (id int, err error)
	CreateSnetOrg(organization SnetOrganization) (id int, err error)
	CreateSnetOrgGroups(orgID int, groups []SnetOrgGroup) (err error)
	GetSnetOrgGroup(groupID string) (SnetOrgGroup, error)
	Health() map[string]string
}

type SnetOrganization struct {
	ID               int        `db:"id"`
	SnetID           string     `db:"snet_id"`
	Name             string     `db:"name"`
	Type             string     `db:"type"`
	Description      string     `db:"description"`
	ShortDescription string     `db:"short_description"`
	URL              string     `db:"url"`
	Owner            string     `db:"owner"`
	Image            string     `db:"image"`
	CreatedAt        time.Time  `db:"created_at"` // not null
	UpdatedAt        time.Time  `db:"updated_at"` // not null
	DeletedAt        *time.Time `db:"deleted_at"` // can be null
}

type SnetService struct {
	ID                    int        `db:"id"`
	SnetID                string     `db:"snet_id"`
	SnetOrgID             string     `db:"snet_org_id"`
	OrgID                 int        `db:"org_id"`
	Version               int        `db:"version"`
	DisplayName           string     `db:"displayname"`
	Encoding              string     `db:"encoding"`
	ServiceType           string     `db:"service_type"`
	ModelIpfsHash         string     `db:"model_ipfs_hash"`
	MPEAddress            string     `db:"mpe_address"`
	URL                   string     `db:"url"`
	Price                 int        `db:"price"`
	GroupID               string     `db:"group_id"`
	FreeCalls             int        `db:"free_calls"`
	FreeCallSignerAddress string     `db:"free_call_signer_address"`
	ShortDescription      string     `db:"short_description"`
	Description           string     `db:"description"`
	CreatedAt             time.Time  `db:"created_at"` // not null
	UpdatedAt             time.Time  `db:"updated_at"` // not null
	DeletedAt             *time.Time `db:"deleted_at"` // can be null
}

type SnetOrgGroup struct {
	ID                         int        `db:"id"`
	OrgID                      int        `db:"org_id"`
	GroupID                    string     `db:"group_id"`
	GroupName                  string     `db:"group_name"`
	PaymentAddress             string     `db:"payment_address"`
	PaymentExpirationThreshold *big.Int   `db:"payment_expiration_threshold"`
	CreatedAt                  time.Time  `db:"created_at"` // not null
	UpdatedAt                  time.Time  `db:"updated_at"` // not null
	DeletedAt                  *time.Time `db:"deleted_at"` // can be null
}
