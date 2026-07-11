CREATE TABLE scans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manifest_digest VARCHAR(255) NOT NULL,
    scanner VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    result_json JSONB
);

CREATE TABLE sboms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manifest_digest VARCHAR(255) NOT NULL,
    format VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    content_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scan_id UUID REFERENCES scans(id) ON DELETE CASCADE,
    cve_id VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    package_name VARCHAR(255) NOT NULL,
    package_version VARCHAR(255) NOT NULL,
    cvss_score NUMERIC(4,2),
    description TEXT
);

CREATE TABLE trust_scores (
    manifest_digest VARCHAR(255) PRIMARY KEY,
    score NUMERIC(5,2) NOT NULL,
    cve_score NUMERIC(5,2) NOT NULL,
    sbom_score NUMERIC(5,2) NOT NULL,
    slsa_score NUMERIC(5,2) NOT NULL,
    signing_score NUMERIC(5,2) NOT NULL,
    freshness_score NUMERIC(5,2) NOT NULL,
    provenance_score NUMERIC(5,2) NOT NULL,
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
