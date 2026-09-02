CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE applications (
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
    ON applications(status);

CREATE INDEX idx_application_created_at
    ON applications(created_at DESC);

CREATE INDEX idx_application_status_created_at
    ON applications(status, created_at DESC);