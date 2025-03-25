-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.requests
    ADD COLUMN file_number VARCHAR(50);

CREATE INDEX idx_requests_file_number ON public.requests(file_number);

ALTER TABLE public.documents
    ADD COLUMN request_id bigint;

ALTER TABLE public.documents
    ADD CONSTRAINT documents_request_id_fkey 
    FOREIGN KEY (request_id)
    REFERENCES public.requests (id) 
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;

CREATE INDEX idx_documents_request_id ON public.documents(request_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.documents
    DROP CONSTRAINT IF EXISTS documents_request_id_fkey;

DROP INDEX IF EXISTS idx_documents_request_id;

ALTER TABLE public.documents
    DROP COLUMN IF EXISTS request_id;

DROP INDEX IF EXISTS idx_requests_file_number;

ALTER TABLE public.requests
    DROP COLUMN IF EXISTS file_number;
-- +goose StatementEnd