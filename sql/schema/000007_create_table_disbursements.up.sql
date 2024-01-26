CREATE TABLE disbursements (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    disbursement_account_guid UUID REFERENCES disbursement_accounts(guid) NOT NULL,
    payment_provider_guid UUID REFERENCES payment_providers(guid) NOT NULL,
    account_name VARCHAR(30) NOT NULL,
    account_number VARCHAR(30) NOT NULL,
    admin_fee DECIMAL CHECK (admin_fee >= 0) NOT NULL DEFAULT 0,
    amount DECIMAL CHECK (amount >= 0) NOT NULL DEFAULT 0,
    amount_with_fee DECIMAL CHECK (amount_with_fee >= 0) NOT NULL DEFAULT 0,
    status disbursement_status NOT NULL DEFAULT 'pending',

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);