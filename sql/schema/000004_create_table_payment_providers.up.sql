CREATE TABLE payment_providers (
    guid UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    admin_fee DECIMAL CHECK (admin_fee >= 0) NOT NULL DEFAULT 0,
    type payment_provider_type NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);