package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"math/big"
	"matrix-ai-framework/internal/config"
	"strings"
	"time"
)

type postgres struct {
	*pgxpool.Pool
}

// New initializes a new postgres connection.
func New() Service {
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.Postgres.User, config.Postgres.Password, config.Postgres.Host, config.Postgres.Port, config.Postgres.Name)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to postgres")
	}
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to take conn from pool")
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `
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
				expires_at          		TIMESTAMP NOT NULL DEFAULT current_timestamp + interval '10 minutes'
			);

	CREATE EXTENSION IF NOT EXISTS pgcrypto;
`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create tables")
		return nil
	}
	db := &postgres{pool}
	return db
}

func (p *postgres) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.Ping(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("db down!")
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

// CreateSnetService creates snet service
func (p *postgres) CreateSnetService(s SnetService) (id int, err error) {
	row := p.Pool.QueryRow(context.Background(),
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
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Can't add snet-service")
	}
	return
}

// CreateSnetOrg creates snet organization
func (p *postgres) CreateSnetOrg(org SnetOrganization) (id int, err error) {
	row := p.Pool.QueryRow(context.Background(),
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
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Can't add snet-org")
	}
	return
}

// CreateSnetOrgGroups creates snet organization group
func (p *postgres) CreateSnetOrgGroups(orgID int, groups []SnetOrgGroup) (err error) {
	ctx := context.Background()
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err)
		return
	}
	defer tx.Rollback(ctx)

	stmt := `
		INSERT INTO snet_org_groups (org_id, group_id, group_name, payment_address, payment_expiration_threshold)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (group_id) DO NOTHING
	`

	for _, group := range groups {
		_, err = tx.Exec(ctx, stmt, orgID, group.GroupID, group.GroupName, group.PaymentAddress, group.PaymentExpirationThreshold)
		if err != nil {
			log.Error().Err(err).Msg("Can't add snet org group")
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetSnetOrgGroup retrieves a snet organization group
func (p *postgres) GetSnetOrgGroup(groupID string) (g SnetOrgGroup, err error) {

	row := p.Pool.QueryRow(context.Background(), "SELECT * FROM snet_org_groups WHERE group_id=$1 AND deleted_at is NULL", groupID)
	var paymentExpirationThreshold int64
	err = row.Scan(&g.ID, &g.OrgID, &g.GroupID, &g.GroupName, &g.PaymentAddress, &paymentExpirationThreshold, &g.CreatedAt, &g.UpdatedAt, &g.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("No org group found with given ID")
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve org group")
		return
	}

	// Convert int64 to *big.Int
	g.PaymentExpirationThreshold = big.NewInt(paymentExpirationThreshold)

	log.Debug().Msgf("Retrieved org group: %+v", g)
	return
}

// GetSnetServices retrieves a list of services
func (p *postgres) GetSnetServices() (services []SnetService, err error) {
	rows, err := p.Pool.Query(context.Background(), "SELECT * FROM snet_services")
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve snet services")
		return services, err
	}
	services, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[SnetService])
	if err != nil {
		log.Error().Err(err).Msg("Failed to scan snet services")
	}
	return services, err
}

// GetSnetOrgs retrieves a list of organizations
func (p *postgres) GetSnetOrgs() ([]SnetOrganization, error) {
	rows, err := p.Pool.Query(context.Background(), "SELECT * FROM snet_organizations WHERE deleted_at is NULL")
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve snet orgs")
		return nil, nil
	}
	orgs, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[SnetOrganization])
	if err != nil {
		log.Error().Err(err)
	}
	return orgs, err
}

// GetSnetService retrieves a snet service
func (p *postgres) GetSnetService(snetID string) (s SnetService, err error) {
	row := p.Pool.QueryRow(context.Background(), "SELECT * FROM snet_services WHERE snet_id=$1 AND deleted_at is NULL", snetID)
	err = row.Scan(&s.ID, &s.SnetID, &s.SnetOrgID, &s.OrgID, &s.Version, &s.DisplayName, &s.Encoding, &s.ServiceType, &s.ModelIpfsHash, &s.MPEAddress, &s.URL, &s.Price, &s.GroupID, &s.FreeCalls, &s.FreeCallSignerAddress, &s.ShortDescription, &s.Description, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("No snet service found with given ID")
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve snet service")
		return
	}
	log.Debug().Msgf("Retrieved snet service: %+v", s)
	return
}

// CreatePaymentState creates a payment state
func (p *postgres) CreatePaymentState(ps *PaymentState) (id uuid.UUID, err error) {
	query := `INSERT INTO payment_states (id, url, key, token_address, to_address, amount) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, ps.ID, ps.URL, ps.Key, ps.TokenAddress, ps.ToAddress, ps.Amount)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment")
	}
	return
}

// GetPaymentStateByKey retrieves a payment state
func (p *postgres) GetPaymentStateByKey(key string) (ps *PaymentState, err error) {
	query := `SELECT * FROM payment_states WHERE key = $1 AND status != 'paid'`
	row := p.Pool.QueryRow(context.Background(), query, key)
	ps = &PaymentState{}
	err = row.Scan(&ps.ID, &ps.URL, &ps.Status, &ps.Key, &ps.TxHash, &ps.TokenAddress, &ps.ToAddress, &ps.Amount, &ps.CreatedAt, &ps.UpdatedAt, &ps.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No payment state found with key %s", key)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve payment state")
		return
	}
	log.Debug().Msgf("Retrieved payment state: %+v", ps)
	return
}

// GetPaymentState retrieves a payment state
func (p *postgres) GetPaymentState(id uuid.UUID) (ps *PaymentState, err error) {
	query := `SELECT * FROM payment_states WHERE id = $1 AND status != 'paid'`
	row := p.Pool.QueryRow(context.Background(), query, id)
	ps = &PaymentState{}
	err = row.Scan(&ps.ID, &ps.URL, &ps.Status, &ps.Key, &ps.TxHash, &ps.TokenAddress, &ps.ToAddress, &ps.Amount, &ps.CreatedAt, &ps.UpdatedAt, &ps.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No payment state found with id %s", id)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve payment state")
		return
	}
	log.Debug().Msgf("Retrieved payment state: %+v", ps)
	return
}

// PatchUpdatePaymentState updates only the specified fields of a payment state
func (p *postgres) PatchUpdatePaymentState(ps *PaymentState) (err error) {
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
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", paramID)
	params = append(params, ps.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
	}

	return
}
