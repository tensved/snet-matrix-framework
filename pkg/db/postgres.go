package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"math/big"
	"strings"
	"time"
)

// postgres wraps the pgxpool.Pool to provide additional methods.
type postgres struct {
	*pgxpool.Pool // Connection pool for PostgreSQL.
}

// New initializes a new PostgreSQL connection pool and returns a Service interface.
//
// Returns:
//   - A Service interface that wraps the PostgreSQL connection pool.
func New() Service {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", config.Postgres.User, config.Postgres.Password, config.Postgres.Host, config.Postgres.Port, config.Postgres.Name)

	pgConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse postgres config")
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to postgres")
		return nil
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to take conn from pool")
		return nil
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS snet_organizations
		(
			id                  SERIAL PRIMARY KEY,
			snet_id 			TEXT NOT NULL UNIQUE,
			name            	TEXT NOT NULL DEFAULT '',
			type 				TEXT NOT NULL DEFAULT '',
			description 		TEXT NOT NULL DEFAULT '',
			short_description 	TEXT NOT NULL DEFAULT '',
			url 				TEXT NOT NULL DEFAULT '',
			owner 				TEXT,
			image 				TEXT,
			created_at          TIMESTAMP NOT NULL DEFAULT current_timestamp,
			updated_at          TIMESTAMP NOT NULL DEFAULT current_timestamp,
			deleted_at          TIMESTAMP DEFAULT NULL
		);

	CREATE TABLE IF NOT EXISTS snet_org_groups
		(
			id                  SERIAL PRIMARY KEY,
			org_id 			INTEGER REFERENCES snet_organizations (id),
			group_id 				TEXT NOT NULL UNIQUE,
			group_name            	TEXT NOT NULL DEFAULT '',
			payment_address 		TEXT NOT NULL DEFAULT '',
			payment_expiration_threshold BIGINT NOT NULL DEFAULT 0,
			created_at          TIMESTAMP NOT NULL DEFAULT current_timestamp,
			updated_at          TIMESTAMP NOT NULL DEFAULT current_timestamp,
			deleted_at          TIMESTAMP DEFAULT null
		);

	CREATE TABLE IF NOT EXISTS snet_services
		(
			id                  		SERIAL PRIMARY KEY,
			snet_id 					TEXT NOT NULL UNIQUE,
			snet_org_id 				TEXT,
			org_id    					INTEGER REFERENCES snet_organizations (id),
			version 					INTEGER,
			displayname         		TEXT NOT NULL DEFAULT '',
			encoding     				TEXT DEFAULT '',
			service_type        		TEXT NOT NULL DEFAULT '',
			model_ipfs_hash     		TEXT NOT NULL DEFAULT '',
			mpe_address 				TEXT NOT NULL DEFAULT '',
			url 						TEXT NOT NULL DEFAULT '',
			price 						INTEGER not null DEFAULT 0,
			group_id 					TEXT NOT NULL DEFAULT '',
			free_calls 					INTEGER NOT NULL DEFAULT 0,
			free_call_signer_address 	TEXT NOT NULL DEFAULT '',
			short_description 			TEXT NOT NULL DEFAULT '',
			description 				TEXT NOT NULL DEFAULT '',
			created_at         			TIMESTAMP NOT NULL DEFAULT current_timestamp,
			updated_at          		TIMESTAMP NOT NULL DEFAULT current_timestamp,
			deleted_at          		TIMESTAMP DEFAULT null
		);

	CREATE TABLE IF NOT EXISTS payment_states
			(
				id                  		UUID PRIMARY KEY,
				url 						TEXT,
				status 						TEXT DEFAULT 'pending',
				key 						TEXT,
				tx_hash 					TEXT DEFAULT NULL,
				token_address 				TEXT NOT NULL,
				to_address 					TEXT NOT NULL,
				amount 						INTEGER NOT NULL DEFAULT 0,
				created_at         			TIMESTAMP NOT NULL DEFAULT current_timestamp,
				updated_at          		TIMESTAMP NOT NULL DEFAULT current_timestamp,
				expires_at          		TIMESTAMP NOT NULL DEFAULT current_timestamp + interval '2 minutes'
			);

	CREATE EXTENSION IF NOT EXISTS pgcrypto;
`)
	if err != nil {
		log.Error().Err(err).Msg("failed to create tables")
		return nil
	}

	return &postgres{pool}
}

// Health checks the health of the database connection.
//
// Returns:
//   - A map with a health status message.
func (p *postgres) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Msg("db down")
		return map[string]string{
			"error": "db down",
		}
	}

	return map[string]string{
		"message": "it's healthy",
	}
}

// CreateSnetService creates a new snet service in the database.
//
// Parameters:
//   - s: An instance of SnetService containing service details.
//
// Returns:
//   - id: The Id of the created service.
//   - error: An error if the operation fails.
func (p *postgres) CreateSnetService(s SnetService) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := p.Pool.QueryRow(ctx,
		`
			INSERT INTO snet_services
   			(snet_id, snet_org_id, org_id, version, displayname, encoding , service_type, model_ipfs_hash, mpe_address, url, price, group_id, free_calls, free_call_signer_address, short_description, description) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			ON CONFLICT (snet_id)
			DO UPDATE SET
			    snet_id=EXCLUDED.snet_id,
				snet_org_id=EXCLUDED.snet_org_id,
				org_id=EXCLUDED.org_id,
				version=EXCLUDED.version,
				displayname=EXCLUDED.displayname,
				encoding=EXCLUDED.encoding,
				service_type=EXCLUDED.service_type,
				model_ipfs_hash=EXCLUDED.model_ipfs_hash,
				mpe_address=EXCLUDED.mpe_address,
				url=EXCLUDED.url,
				price=EXCLUDED.price,
				group_id=EXCLUDED.group_id,
				free_calls=EXCLUDED.free_calls,
				free_call_signer_address=EXCLUDED.free_call_signer_address,
				short_description=EXCLUDED.short_description,
				description=EXCLUDED.description
			RETURNING id`,
		s.SnetID, s.SnetOrgID, s.OrgID, s.Version, s.DisplayName, s.Encoding, s.ServiceType, s.ModelIpfsHash, s.MPEAddress, s.URL, s.Price, s.GroupID, s.FreeCalls, s.FreeCallSignerAddress, s.ShortDescription, s.Description)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, errors.New("failed to add snet-service")
	}
	return id, nil
}

// CreateSnetOrg creates a new snet organization in the database.
//
// Parameters:
//   - org: An instance of SnetOrganization containing organization details.
//
// Returns:
//   - id: The Id of the created organization.
//   - error: An error if the operation fails.
func (p *postgres) CreateSnetOrg(org SnetOrganization) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := p.Pool.QueryRow(ctx,
		`
			INSERT INTO snet_organizations
    		(snet_id, name, type, short_description, description, url, owner, image)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (snet_id)
			DO UPDATE SET
			    snet_id=EXCLUDED.snet_id,
			    name=EXCLUDED.name,
			    type=EXCLUDED.type,
			    short_description=EXCLUDED.short_description,
			    description=EXCLUDED.description,
			    url=EXCLUDED.url,
			    owner=EXCLUDED.owner,
			    image=EXCLUDED.image
			RETURNING id`,
		org.SnetID, org.Name, org.Type, org.ShortDescription, org.Description, org.URL, org.Owner, org.Image)
	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("failed to add snet org")
		return 0, err
	}
	return id, nil
}

// CreateSnetOrgGroups creates multiple snet organization groups in the database.
//
// Parameters:
//   - orgID: The Id of the organization to which the groups belong.
//   - groups: A slice of SnetOrgGroup containing group details.
//
// Returns:
//   - error: An error if the operation fails.
func (p *postgres) CreateSnetOrgGroups(orgID int, groups []SnetOrgGroup) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Error().Err(err)
		}
	}(tx, ctx)

	stmt := `
		INSERT INTO snet_org_groups (org_id, group_id, group_name, payment_address, payment_expiration_threshold)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (group_id) DO NOTHING
	`

	for _, group := range groups {
		_, err = tx.Exec(ctx, stmt, orgID, group.GroupID, group.GroupName, group.PaymentAddress, group.PaymentExpirationThreshold)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetSnetOrgGroup retrieves a snet organization group from the database.
//
// Parameters:
//   - groupID: The Id of the group to retrieve.
//
// Returns:
//   - g: The retrieved SnetOrgGroup instance.
//   - error: An error if the operation fails.
func (p *postgres) GetSnetOrgGroup(groupID string) (SnetOrgGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := p.Pool.QueryRow(ctx, "SELECT * FROM snet_org_groups WHERE group_id=$1 AND deleted_at is NULL", groupID)
	var paymentExpirationThreshold int64
	g := SnetOrgGroup{}
	err := row.Scan(&g.ID, &g.OrgID, &g.GroupID, &g.GroupName, &g.PaymentAddress, &paymentExpirationThreshold, &g.CreatedAt, &g.UpdatedAt, &g.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SnetOrgGroup{}, errors.New("snet org not found")
		}
		return SnetOrgGroup{}, err
	}

	// Convert int64 to *big.Int
	g.PaymentExpirationThreshold = big.NewInt(paymentExpirationThreshold)

	log.Debug().Msgf("retrieved org group: %+v", g)
	return g, nil
}

// GetSnetServices retrieves a list of snet services from the database.
//
// Returns:
//   - services: A slice of SnetService instances.
//   - error: An error if the operation fails.
func (p *postgres) GetSnetServices() ([]SnetService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := p.Pool.Query(ctx, "SELECT * FROM snet_services")
	if err != nil {
		return nil, err
	}
	services, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[SnetService])
	if err != nil {
		return nil, err
	}
	return services, nil
}

// GetSnetOrgs retrieves a list of snet organizations from the database.
//
// Returns:
//   - orgs: A slice of SnetOrganization instances.
//   - error: An error if the operation fails.
func (p *postgres) GetSnetOrgs() ([]SnetOrganization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := p.Pool.Query(ctx, "SELECT * FROM snet_organizations WHERE deleted_at is NULL")
	if err != nil {
		return nil, errors.New("failed to retrieve snet orgs")
	}
	orgs, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[SnetOrganization])
	if err != nil {
		return nil, errors.New("failed to scan snet orgs")
	}
	return orgs, nil
}

// GetSnetService retrieves a snet service from the database.
//
// Parameters:
//   - snetID: The Id of the service to retrieve.
//
// Returns:
//   - s: The retrieved SnetService instance.
//   - error: An error if the operation fails.
func (p *postgres) GetSnetService(snetID string) (*SnetService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := &SnetService{}
	row := p.Pool.QueryRow(ctx, "SELECT * FROM snet_services WHERE snet_id=$1 AND deleted_at is NULL", snetID)
	err := row.Scan(&s.ID, &s.SnetID, &s.SnetOrgID, &s.OrgID, &s.Version, &s.DisplayName, &s.Encoding, &s.ServiceType, &s.ModelIpfsHash, &s.MPEAddress, &s.URL, &s.Price, &s.GroupID, &s.FreeCalls, &s.FreeCallSignerAddress, &s.ShortDescription, &s.Description, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("no snet service found with given id")
		}
		return nil, err
	}
	log.Debug().Msgf("retrieved snet service: %+v", s)
	return s, nil
}

// CreatePaymentState creates a new payment state in the database.
//
// Parameters:
//   - ps: An instance of PaymentState containing payment details.
//
// Returns:
//   - id: The UUID of the created payment state.
//   - error: An error if the operation fails.
func (p *postgres) CreatePaymentState(ps *PaymentState) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `INSERT INTO payment_states (id, url, key, token_address, to_address, amount) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := p.Pool.QueryRow(ctx, query, ps.ID, ps.URL, ps.Key, ps.TokenAddress, ps.ToAddress, ps.Amount)
	id := uuid.UUID{}
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, errors.New("failed to create payment")
	}
	return id, nil
}

// GetPaymentStateByKey retrieves a payment state from the database using a key.
//
// Parameters:
//   - key: The key to retrieve the payment state.
//
// Returns:
//   - ps: The retrieved PaymentState instance.
//   - error: An error if the operation fails.
func (p *postgres) GetPaymentStateByKey(key string) (*PaymentState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT * FROM payment_states WHERE key = $1 AND status != 'paid'`
	row := p.Pool.QueryRow(ctx, query, key)
	ps := &PaymentState{}
	err := row.Scan(&ps.ID, &ps.URL, &ps.Status, &ps.Key, &ps.TxHash, &ps.TokenAddress, &ps.ToAddress, &ps.Amount, &ps.CreatedAt, &ps.UpdatedAt, &ps.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no payment state found with key %s", key)
		}
		return nil, errors.New("failed to retrieve payment state")
	}
	log.Debug().Msgf("retrieved payment state: %+v", ps)
	return ps, nil
}

// GetPaymentState retrieves a payment state from the database using an Id.
//
// Parameters:
//   - id: The UUID of the payment state to retrieve.
//
// Returns:
//   - ps: The retrieved PaymentState instance.
//   - error: An error if the operation fails.
func (p *postgres) GetPaymentState(id uuid.UUID) (*PaymentState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT * FROM payment_states WHERE id = $1 AND status != 'paid'`
	row := p.Pool.QueryRow(ctx, query, id)
	ps := &PaymentState{}
	err := row.Scan(&ps.ID, &ps.URL, &ps.Status, &ps.Key, &ps.TxHash, &ps.TokenAddress, &ps.ToAddress, &ps.Amount, &ps.CreatedAt, &ps.UpdatedAt, &ps.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("no payment state found")
		}
		return nil, err
	}
	log.Debug().Msgf("retrieved payment state: %+v", ps)
	return ps, nil
}

// PatchUpdatePaymentState updates specific fields of a payment state in the database.
//
// Parameters:
//   - ps: An instance of PaymentState containing the fields to update.
//
// Returns:
//   - error: An error if the operation fails.
func (p *postgres) PatchUpdatePaymentState(ps *PaymentState) error {
	query := "UPDATE payment_states SET updated_at = NOW(), "
	updates := []string{}
	params := []any{}
	paramID := 1

	if ps.Status != "" {
		updates = append(updates, fmt.Sprintf("status = $%d", paramID))
		params = append(params, ps.Status)
		paramID++
	}
	if ps.TxHash != nil {
		updates = append(updates, fmt.Sprintf("tx_hash = $%d", paramID))
		params = append(params, ps.TxHash)
		paramID++
	}

	if len(updates) == 0 {
		return errors.New("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", paramID)
	params = append(params, ps.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := p.Pool.Exec(ctx, query, params...)
	if err != nil {
		return errors.New("failed to update payment state")
	}

	return nil
}
