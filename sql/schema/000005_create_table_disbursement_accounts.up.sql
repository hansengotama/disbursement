CREATE TABLE disbursement_accounts (
    guid UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    payment_provider_guid UUID REFERENCES payment_providers(guid) NOT NULL,
    name VARCHAR(30) NOT NULL,
    number VARCHAR(30) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);