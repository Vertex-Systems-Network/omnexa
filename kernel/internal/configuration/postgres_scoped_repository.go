package configuration

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresScopedRepository persists only kernel.configuration-owned scoped overrides.
type PostgresScopedRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresScopedRepository(pool *pgxpool.Pool) (*PostgresScopedRepository, error) {
	if pool == nil {
		return nil, scopedServiceInvalidFailure()
	}
	return &PostgresScopedRepository{pool: pool}, nil
}

func (repository *PostgresScopedRepository) resolveExact(
	ctx context.Context,
	key Key,
	tenantID tenancy.TenantID,
	organizationID organization.NodeID,
) (ProviderResult, error) {
	if repository == nil || repository.pool == nil || !keyPattern.MatchString(string(key)) || !tenantID.Valid() || (organizationID != "" && !organizationID.Valid()) {
		return ProviderResult{}, scopedContextInvalidFailure()
	}

	var kind string
	var encoded string
	var revision uint64
	queryErr := repository.pool.QueryRow(
		ctx,
		`SELECT value_kind, value_text, revision
         FROM omnexa_configuration.setting_overrides
         WHERE tenant_id = $1
           AND organization_id IS NOT DISTINCT FROM $2::uuid
           AND setting_key = $3`,
		string(tenantID), nullableOrganizationID(organizationID), string(key),
	).Scan(&kind, &encoded, &revision)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return ProviderResult{}, ErrProviderValueNotFound
	}
	if queryErr != nil {
		return ProviderResult{}, scopedPersistenceFailure(queryErr)
	}
	value, err := decodeScopedValue(Kind(kind), encoded)
	if err != nil || revision == 0 {
		return ProviderResult{}, scopedPersistenceFailure(errors.New("stored scoped setting is invalid"))
	}
	return ProviderResult{Value: value, Revision: revision}, nil
}

func (repository *PostgresScopedRepository) upsert(
	ctx context.Context,
	key Key,
	value Value,
	tenantID tenancy.TenantID,
	organizationID organization.NodeID,
	changedAt time.Time,
) (uint64, error) {
	if repository == nil || repository.pool == nil || !keyPattern.MatchString(string(key)) || !value.valid() || !tenantID.Valid() ||
		(organizationID != "" && !organizationID.Valid()) || changedAt.IsZero() {
		return 0, scopedValueInvalidFailure()
	}
	encoded, err := encodeScopedValue(value)
	if err != nil {
		return 0, err
	}
	instant := changedAt.UTC()
	var revision uint64
	if organizationID == "" {
		err = repository.pool.QueryRow(
			ctx,
			`INSERT INTO omnexa_configuration.setting_overrides
                (tenant_id, organization_id, setting_key, value_kind, value_text, revision, created_at, updated_at)
             VALUES ($1, NULL, $2, $3, $4, 1, $5, $5)
             ON CONFLICT (tenant_id, setting_key) WHERE organization_id IS NULL
             DO UPDATE SET
                value_kind = EXCLUDED.value_kind,
                value_text = EXCLUDED.value_text,
                revision = omnexa_configuration.setting_overrides.revision + 1,
                updated_at = EXCLUDED.updated_at
             RETURNING revision`,
			string(tenantID), string(key), string(value.Kind()), encoded, instant,
		).Scan(&revision)
	} else {
		err = repository.pool.QueryRow(
			ctx,
			`INSERT INTO omnexa_configuration.setting_overrides
                (tenant_id, organization_id, setting_key, value_kind, value_text, revision, created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5, 1, $6, $6)
             ON CONFLICT (tenant_id, organization_id, setting_key) WHERE organization_id IS NOT NULL
             DO UPDATE SET
                value_kind = EXCLUDED.value_kind,
                value_text = EXCLUDED.value_text,
                revision = omnexa_configuration.setting_overrides.revision + 1,
                updated_at = EXCLUDED.updated_at
             RETURNING revision`,
			string(tenantID), string(organizationID), string(key), string(value.Kind()), encoded, instant,
		).Scan(&revision)
	}
	if err != nil {
		return 0, scopedPersistenceFailure(err)
	}
	return revision, nil
}

func nullableOrganizationID(id organization.NodeID) any {
	if id == "" {
		return nil
	}
	return string(id)
}

func encodeScopedValue(value Value) (string, error) {
	if !value.valid() {
		return "", scopedValueInvalidFailure()
	}
	switch value.Kind() {
	case KindBool:
		actual, _ := value.Bool()
		return strconv.FormatBool(actual), nil
	case KindString:
		actual, _ := value.String()
		return actual, nil
	case KindInt:
		actual, _ := value.Int()
		return strconv.FormatInt(actual, 10), nil
	case KindDuration:
		actual, _ := value.Duration()
		return strconv.FormatInt(int64(actual), 10), nil
	default:
		return "", scopedValueInvalidFailure()
	}
}

func decodeScopedValue(kind Kind, encoded string) (Value, error) {
	switch kind {
	case KindBool:
		value, err := strconv.ParseBool(encoded)
		if err != nil {
			return Value{}, scopedValueInvalidFailure()
		}
		return BoolValue(value), nil
	case KindString:
		value := StringValue(encoded)
		if !value.valid() {
			return Value{}, scopedValueInvalidFailure()
		}
		return value, nil
	case KindInt:
		value, err := strconv.ParseInt(encoded, 10, 64)
		if err != nil {
			return Value{}, scopedValueInvalidFailure()
		}
		return IntValue(value), nil
	case KindDuration:
		value, err := strconv.ParseInt(encoded, 10, 64)
		if err != nil {
			return Value{}, scopedValueInvalidFailure()
		}
		return DurationValue(time.Duration(value)), nil
	default:
		return Value{}, scopedValueInvalidFailure()
	}
}
