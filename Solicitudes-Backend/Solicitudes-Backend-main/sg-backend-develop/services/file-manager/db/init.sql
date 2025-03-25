-- Document types table
CREATE TABLE IF NOT EXISTS document_types (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_mandatory BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Documents table
CREATE TABLE IF NOT EXISTS documents (
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

-- Create trigger function for updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for tables that need updated_at
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
