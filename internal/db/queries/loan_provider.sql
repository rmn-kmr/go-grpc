
-- name: GetLoanProviders :many
SELECT * FROM loan_provider;

-- name: GetLoanProviderBySlug :one
SELECT * FROM loan_provider
WHERE provider = $1 LIMIT 1;

-- name: GetLoanProviderByID :one
SELECT * FROM loan_provider
WHERE provider_id = $1 LIMIT 1;