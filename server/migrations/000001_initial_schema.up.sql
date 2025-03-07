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
   latitude TEXT NOT NULL,
   longitude TEXT NOT NULL,
   info TEXT NOT NULL,
    street TEXT NOT NULL,
    number INT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    formatted_address TEXT NOT NULL,
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
    lead BOOLEAN DEFAULT false,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE RESTRICT,
    UNIQUE (project_id, account_id)
);

INSERT INTO accounts (first, last, password, email, created_at, updated_at)
VALUES ('admin', 'user', 'password', 'scarrington@gmail.com', LOCALTIMESTAMP, LOCALTIMESTAMP);

INSERT INTO locations (latitude, longitude, info, street, number, city, state, postal_code, formatted_address, created_at, updated_at)
VALUES ('39.5023509486161', '-104.87569755087917', 'Journey Church', 'Clydesdale Road', '9009', 'Castle Rock', 'Colorado', '80108', '9009 Clydesdale Rd, Castle Rock, CO 80108', LOCALTIMESTAMP, LOCALTIMESTAMP);

INSERT INTO projects (google_id, name, required, needed, admin_id, location_id, created_at, updated_at)
VALUES ('google1', 'Base Project Example', 200, 159, 1, 1, LOCALTIMESTAMP, LOCALTIMESTAMP);
