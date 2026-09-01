CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE application (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                       name TEXT NOT NULL,
                       phone TEXT NOT NULL,
                       comment TEXT NOT NULL DEFAULT '',
                       source TEXT NOT NULL DEFAULT '',

                       status TEXT NOT NULL DEFAULT 'new'
                           CHECK (
                               status IN (
                                          'new',
                                          'in_progress',
                                          'success',
                                          'rejected'
                                   )
                               ),

                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_application_status
    ON application(status);

CREATE INDEX idx_application_created_at
    ON application(created_at DESC);

CREATE INDEX idx_application_status_created_at
    ON application(status, created_at DESC);