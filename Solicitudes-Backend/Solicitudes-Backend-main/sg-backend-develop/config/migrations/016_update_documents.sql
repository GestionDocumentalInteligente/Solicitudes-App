-- 014_update_document_tables.sql
-- +goose Up
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS document_types;

CREATE TABLE document_types (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  is_mandatory BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE documents (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  code VARCHAR(255) NOT NULL, 
  document_type_id INTEGER REFERENCES document_types(id),
  file_id VARCHAR(255) NOT NULL,  
  file_url VARCHAR(255),
  is_verified BOOLEAN DEFAULT false,
  status BOOLEAN DEFAULT false,
  content TEXT NOT NULL,
  observations TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS document_types;