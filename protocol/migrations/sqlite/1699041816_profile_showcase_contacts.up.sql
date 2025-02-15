DROP TABLE profile_showcase_preferences;
DROP TABLE profile_showcase_contacts;

-- Four tables for storing own profile showcase preferences
CREATE TABLE profile_showcase_communities_preferences (
    community_id TEXT PRIMARY KEY ON CONFLICT REPLACE,
    visibility INT NOT NULL DEFAULT 0,
    sort_order INT DEFAULT 0
);

CREATE TABLE profile_showcase_accounts_preferences (
    address TEXT PRIMARY KEY ON CONFLICT REPLACE,
    name TEXT DEFAULT "",
    color_id DEFAULT "",
    emoji VARCHAR DEFAULT "",
    visibility INT NOT NULL DEFAULT 0,
    sort_order INT DEFAULT 0
);

CREATE TABLE profile_showcase_collectibles_preferences (
    uid TEXT PRIMARY KEY ON CONFLICT REPLACE,
    visibility INT NOT NULL DEFAULT 0,
    sort_order INT DEFAULT 0
);

CREATE TABLE profile_showcase_assets_preferences (
    symbol TEXT PRIMARY KEY ON CONFLICT REPLACE,
    visibility INT NOT NULL DEFAULT 0,
    sort_order INT DEFAULT 0
);

-- Four tables for storing profile showcase for each contact
CREATE TABLE profile_showcase_communities_contacts (
    community_id TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    contact_id TEXT NOT NULL,
    PRIMARY KEY (contact_id, community_id)
);
CREATE INDEX profile_showcase_communities_contact_id ON profile_showcase_communities_contacts (contact_id);

CREATE TABLE profile_showcase_accounts_contacts (
    address TEXT NOT NULL,
    name TEXT DEFAULT "",
    color_id DEFAULT "",
    emoji VARCHAR DEFAULT "",
    sort_order INT DEFAULT 0,
    contact_id TEXT NOT NULL,
    PRIMARY KEY (contact_id, address)
);
CREATE INDEX profile_showcase_accounts_contact_id ON profile_showcase_accounts_contacts (contact_id);

CREATE TABLE profile_showcase_collectibles_contacts (
    uid TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    contact_id TEXT NOT NULL,
    PRIMARY KEY (contact_id, uid)
);
CREATE INDEX profile_showcase_collectibles_contact_id ON profile_showcase_collectibles_contacts (contact_id);

CREATE TABLE profile_showcase_assets_contacts (
    symbol TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    contact_id TEXT NOT NULL,
    PRIMARY KEY (contact_id, symbol)
);
CREATE INDEX profile_showcase_assets_contact_id ON profile_showcase_assets_contacts (contact_id);
