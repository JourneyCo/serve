CREATE TABLE IF NOT EXISTS accounts(
    id serial PRIMARY KEY,
    first TEXT NOT NULL,
    last TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS locations(
   id serial PRIMARY KEY,
   address TEXT NOT NULL,
   latitude TEXT NOT NULL,
   longitude TEXT NOT NULL,
   info TEXT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL,
   updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS projects(
    id serial PRIMARY KEY,
    google_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    required TEXT NOT NULL,
    needed TEXT NOT NULL,
    admin_id BIGINT NOT NULL,
    location_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (admin_id) REFERENCES accounts(id) ON DELETE RESTRICT,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS registrations(
    project_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE RESTRICT
)