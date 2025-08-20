package db

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for database operations related to Snet organizations, services, and payment states.
type Service interface {
	GetSnetOrgs() ([]SnetOrganization, error)                                // Retrieves a list of Snet organizations.
	GetSnetServices() ([]SnetService, error)                                 // Retrieves a list of Snet services.
	GetSnetService(snetID string) (s *SnetService, err error)                // Retrieves a specific Snet service by its Id.
	CreateSnetService(service SnetService) (id int, err error)               // Creates a new Snet service.
	CreateSnetOrg(organization SnetOrganization) (id int, err error)         // Creates a new Snet organization.
	CreateSnetOrgGroups(orgID int, groups []SnetOrgGroup) (err error)        // Creates multiple Snet organization groups.
	GetSnetOrgGroup(groupID string) (SnetOrgGroup, error)                    // Retrieves a specific Snet organization group by its Id.
	CreatePaymentState(paymentState *PaymentState) (id uuid.UUID, err error) // Creates a new payment state.
	GetPaymentState(id uuid.UUID) (ps *PaymentState, err error)              // Retrieves a specific payment state by its UUID.
	GetPaymentStateByKey(key string) (ps *PaymentState, err error)           // Retrieves a payment state by its key.
	PatchUpdatePaymentState(ps *PaymentState) (err error)                    // Updates specific fields of a payment state.
	Health() map[string]string                                               // Checks the health of the database connection.
}

// SnetOrganization represents an organization in the Snet system.
type SnetOrganization struct {
	ID               int        `db:"id"`                // The ID of the organization.
	SnetID           string     `db:"snet_id"`           // The unique Snet ID of the organization.
	Name             string     `db:"name"`              // The name of the organization.
	Type             string     `db:"type"`              // The type of the organization.
	Description      string     `db:"description"`       // The description of the organization.
	ShortDescription string     `db:"short_description"` // The short description of the organization.
	URL              string     `db:"url"`               // The URL of the organization's website.
	Owner            string     `db:"owner"`             // The owner of the organization.
	Image            string     `db:"image"`             // The image URL of the organization.
	CreatedAt        time.Time  `db:"created_at"`        // The creation timestamp of the organization.
	UpdatedAt        time.Time  `db:"updated_at"`        // The last update timestamp of the organization.
	DeletedAt        *time.Time `db:"deleted_at"`        // The deletion timestamp of the organization, can be null.
}

// SnetService represents a service in the Snet system.
type SnetService struct {
	ID                    int        `db:"id"`                       // The ID of the service.
	SnetID                string     `db:"snet_id"`                  // The unique Snet ID of the service.
	SnetOrgID             string     `db:"snet_org_id"`              // The Snet organization ID associated with the service.
	OrgID                 int        `db:"org_id"`                   // The organization ID associated with the service.
	Version               int        `db:"version"`                  // The version of the service.
	DisplayName           string     `db:"displayname"`              // The display name of the service.
	Encoding              string     `db:"encoding"`                 // The encoding type of the service.
	ServiceType           string     `db:"service_type"`             // The type of the service.
	ModelIpfsHash         string     `db:"model_ipfs_hash"`          // The IPFS hash of the service model.
	ServiceApiSource      string     `db:"service_api_source"`       // The service API source (new field for newer services).
	MPEAddress            string     `db:"mpe_address"`              // The MPE address of the service.
	URL                   string     `db:"url"`                      // The URL of the service.
	Price                 int        `db:"price"`                    // The price of the service.
	GroupID               string     `db:"group_id"`                 // The group ID associated with the service.
	FreeCalls             int        `db:"free_calls"`               // The number of free calls allowed for the service.
	FreeCallSignerAddress string     `db:"free_call_signer_address"` // The address of the free call signer.
	ShortDescription      string     `db:"short_description"`        // The short description of the service.
	Description           string     `db:"description"`              // The description of the service.
	CreatedAt             time.Time  `db:"created_at"`               // The creation timestamp of the service.
	UpdatedAt             time.Time  `db:"updated_at"`               // The last update timestamp of the service.
	DeletedAt             *time.Time `db:"deleted_at"`               // The deletion timestamp of the service, can be null.
}

// SnetOrgGroup represents a group within an organization in the Snet system.
type SnetOrgGroup struct {
	ID                         int        `db:"id"`                           // The ID of the group.
	OrgID                      int        `db:"org_id"`                       // The organization ID associated with the group.
	GroupID                    string     `db:"group_id"`                     // The unique ID of the group.
	GroupName                  string     `db:"group_name"`                   // The name of the group.
	PaymentAddress             string     `db:"payment_address"`              // The payment address associated with the group.
	PaymentExpirationThreshold *big.Int   `db:"payment_expiration_threshold"` // The payment expiration threshold for the group.
	CreatedAt                  time.Time  `db:"created_at"`                   // The creation timestamp of the group.
	UpdatedAt                  time.Time  `db:"updated_at"`                   // The last update timestamp of the group.
	DeletedAt                  *time.Time `db:"deleted_at"`                   // The deletion timestamp of the group, can be null.
}

// PaymentState represents the state of a user interacting with the bot for payments.
type PaymentState struct {
	ID           uuid.UUID `json:"id" db:"id"`                      // The UUID of the payment.
	URL          string    `json:"url" db:"url"`                    // The URL for the transaction request.
	Key          string    `json:"key" db:"key"`                    // The key for the payment state, typically in the format "{roomId} {userId} {serviceName} {methodName}".
	TxHash       *string   `json:"txHash" db:"tx_hash"`             // The transaction hash, can be null.
	TokenAddress string    `json:"tokenAddress" db:"token_address"` // The token address for the payment.
	ToAddress    string    `json:"toAddress" db:"to_address"`       // The recipient address for the payment.
	Amount       int       `json:"amount" db:"amount"`              // The amount for the payment.
	Status       string    `json:"status" db:"status"`              // The status of the payment, e.g., pending, expired, or paid.
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`       // The creation timestamp of the payment state.
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`       // The last update timestamp of the payment state.
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`       // The expiration timestamp of the payment state.
}
