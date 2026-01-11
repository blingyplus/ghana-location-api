-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Countries table
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(4) UNIQUE NOT NULL,
    name VARCHAR NOT NULL
);

-- Regions table
CREATE TABLE regions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    slug VARCHAR UNIQUE NOT NULL,
    capital VARCHAR
);

-- Districts table
CREATE TABLE districts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    slug VARCHAR UNIQUE NOT NULL,
    type VARCHAR NOT NULL CHECK (type IN ('metro', 'municipal', 'district')),
    capital VARCHAR
);

-- Constituencies table
CREATE TABLE constituencies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    district_id UUID REFERENCES districts(id) ON DELETE SET NULL,
    name VARCHAR NOT NULL,
    slug VARCHAR UNIQUE NOT NULL
);

-- Cities table
CREATE TABLE cities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    district_id UUID NOT NULL REFERENCES districts(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    lat DECIMAL(10, 8),
    lng DECIMAL(11, 8)
);

-- Indexes for performance
CREATE INDEX idx_regions_country_id ON regions(country_id);
CREATE INDEX idx_regions_slug ON regions(slug);
CREATE INDEX idx_districts_region_id ON districts(region_id);
CREATE INDEX idx_districts_slug ON districts(slug);
CREATE INDEX idx_constituencies_district_id ON constituencies(district_id);
CREATE INDEX idx_constituencies_slug ON constituencies(slug);
CREATE INDEX idx_cities_district_id ON cities(district_id);
