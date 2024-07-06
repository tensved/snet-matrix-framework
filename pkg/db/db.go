package db

import (
	"github.com/google/uuid"
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
	CreatePaymentState(paymentState *PaymentState) (id uuid.UUID, err error)
	GetPaymentState(id uuid.UUID) (ps *PaymentState, err error)
	GetPaymentStateByKey(key string) (ps *PaymentState, err error)
	PatchUpdatePaymentState(ps *PaymentState) (err error)
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

// PaymentState represents the state of a user interacting with the bot for payments
type PaymentState struct {
	ID           uuid.UUID `json:"id" db:"id"`   // UUID of payment
	URL          string    `json:"url" db:"url"` // URl for transaction request
	Key          string    `json:"key" db:"key"` // key: "{roomId} {userId} {serviceName} {methodName}"
	TxHash       *string   `json:"txHash" db:"tx_hash"`
	TokenAddress string    `json:"tokenAddress" db:"token_address"`
	ToAddress    string    `json:"toAddress" db:"to_address"`
	Amount       int       `json:"amount" db:"amount"`
	Status       string    `json:"status" db:"status"` // status of payment – pending, expired or paid
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
}
