-- Install extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enums for type safety
CREATE TYPE application_stage AS ENUM (
'applied', 
'screening', 
'interview', 
'offer', 
'rejected'
);

CREATE TYPE experience_level AS ENUM (
'entry', 
'mid', 
'senior', 
'lead'
);

-- Tables
CREATE TABLE files (
    id UUID PRIMARY KEY,
    path VARCHAR(255) NOT NULL UNIQUE,
    file_type VARCHAR(50) NOT NULL,
    checksum VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE candidates (
    id UUID PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(255) NOT NULL,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE RESTRICT,
    status VARCHAR(255) NOT NULL DEFAULT 'new' 
        CHECK (status IN ('new', 'reviewed', 'contacted', 'hired', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE skills (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    category VARCHAR(255) NOT NULL 
        CHECK (category IN ('technical', 'soft', 'language', 'certification')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE departments (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE job_roles (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    department_id UUID NOT NULL REFERENCES departments(id),
    level experience_level NOT NULL,
    salary_range VARCHAR(255),
    location VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE job_role_requirements (
    job_role_id UUID NOT NULL REFERENCES job_roles(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    required BOOLEAN NOT NULL DEFAULT true,
    min_experience_years NUMERIC(3,1),
    importance INT NOT NULL CHECK (importance BETWEEN 1 AND 5),
    PRIMARY KEY (job_role_id, skill_id)
);

CREATE TABLE applications (
    id UUID PRIMARY KEY,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    job_role_id UUID NOT NULL REFERENCES job_roles(id) ON DELETE RESTRICT,
    stage application_stage NOT NULL DEFAULT 'applied',
    score NUMERIC(5,2),
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (candidate_id, job_role_id)
);

CREATE TABLE candidate_skills (
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    years_experience NUMERIC(4,1),
    last_used DATE,
    PRIMARY KEY (candidate_id, skill_id)
);

CREATE TABLE screening_results (
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    model_version VARCHAR(255) NOT NULL,
    raw_response JSONB NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (application_id)
);

CREATE TABLE screening_criteria (
    id UUID PRIMARY KEY,
    screening_result_id UUID NOT NULL REFERENCES screening_results(id) ON DELETE CASCADE,
    criteria_text VARCHAR(255) NOT NULL,
    decision BOOLEAN NOT NULL,
    reasoning VARCHAR(255) NOT NULL,
    matched_skills UUID[],
    missing_skills UUID[]
);

-- Indexes for common queries
CREATE INDEX idx_candidates_email ON candidates(email);
CREATE INDEX idx_job_roles_active ON job_roles(is_active);
CREATE INDEX idx_applications_score ON applications(score);
CREATE INDEX idx_skills_category ON skills(category);
CREATE INDEX idx_screening_processed ON screening_results(processed_at);
