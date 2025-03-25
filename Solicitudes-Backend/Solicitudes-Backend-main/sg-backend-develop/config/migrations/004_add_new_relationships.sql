-- +goose Up
-- +goose StatementBegin

-- Addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    street VARCHAR(255) NOT NULL,
    street_number VARCHAR(20) NOT NULL,
    floor VARCHAR(20),
    apartment VARCHAR(20),
    functional_unit VARCHAR(50),
    postal_code VARCHAR(20),
    city VARCHAR(100),
    state VARCHAR(100),
    is_common_area BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Request status table
CREATE TABLE IF NOT EXISTS request_status (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    requires_review BOOLEAN DEFAULT false,
    is_final_state BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Document types table
CREATE TABLE IF NOT EXISTS document_types (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_mandatory BOOLEAN DEFAULT false,
    required_metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Requests table
CREATE TABLE IF NOT EXISTS requests (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    request_type_id INT REFERENCES request_types(id),
    address_id BIGINT REFERENCES addresses(id),
    status_id INT REFERENCES request_status(id),
    description TEXT,
    verification_complete BOOLEAN DEFAULT false,
    verification_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Documents table
CREATE TABLE IF NOT EXISTS documents (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    request_id BIGINT REFERENCES requests(id),
    document_type_id INT REFERENCES document_types(id),
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50),
    file_url VARCHAR(255),
    is_verified BOOLEAN DEFAULT false,
    observations TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    request_id BIGINT REFERENCES requests(id),
    user_id BIGINT REFERENCES users(id),
    content TEXT NOT NULL,
    sent_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workflow history table
CREATE TABLE IF NOT EXISTS workflow_history (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    request_id BIGINT REFERENCES requests(id),
    previous_status_id INT REFERENCES request_status(id),
    new_status_id INT REFERENCES request_status(id),
    user_id BIGINT REFERENCES users(id),
    change_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    observations TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Verifications table
CREATE TABLE IF NOT EXISTS verifications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    request_id BIGINT REFERENCES requests(id),
    user_id BIGINT REFERENCES users(id),
    result BOOLEAN NOT NULL,
    observations TEXT,
    verification_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_requests_user_id ON requests(user_id);
CREATE INDEX IF NOT EXISTS idx_requests_status_id ON requests(status_id);
CREATE INDEX IF NOT EXISTS idx_documents_request_id ON documents(request_id);
CREATE INDEX IF NOT EXISTS idx_workflow_history_request_id ON workflow_history(request_id);
CREATE INDEX IF NOT EXISTS idx_messages_request_id ON messages(request_id);
CREATE INDEX IF NOT EXISTS idx_verifications_request_id ON verifications(request_id);

-- Create trigger function for updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Create triggers for tables that need updated_at
DO $$
BEGIN
   IF NOT EXISTS (
      SELECT 1 
      FROM pg_trigger 
      WHERE tgname = 'update_requests_updated_at'
   ) THEN
      CREATE TRIGGER update_requests_updated_at
      BEFORE UPDATE ON requests
      FOR EACH ROW
      EXECUTE FUNCTION update_updated_at_column();
   END IF;
END $$;

DO $$
BEGIN
   IF NOT EXISTS (
      SELECT 1 
      FROM pg_trigger 
      WHERE tgname = 'update_documents_updated_at'
   ) THEN
      CREATE TRIGGER update_documents_updated_at
      BEFORE UPDATE ON documents
      FOR EACH ROW
      EXECUTE FUNCTION update_updated_at_column();
   END IF;
END $$;

-- Seed data for request_status table
INSERT INTO request_status (name, description, requires_review, is_final_state, created_at)
VALUES
('Pending', 'Request is pending review', true, false, CURRENT_TIMESTAMP),
('Approved', 'Request has been approved', false, true, CURRENT_TIMESTAMP),
('Rejected', 'Request has been rejected', false, true, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS verifications;
DROP TABLE IF EXISTS workflow_history;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS document_types;
DROP TABLE IF EXISTS request_status;
DROP TABLE IF EXISTS addresses;

DROP FUNCTION IF EXISTS update_updated_at_column();

-- +goose StatementEnd
