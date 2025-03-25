CREATE TABLE IF NOT EXISTS accounts(
    id TEXT PRIMARY KEY,
    first TEXT DEFAULT '',
    last TEXT DEFAULT '',
    email TEXT UNIQUE,
    cellphone TEXT DEFAULT '',
    text_permission BOOLEAN DEFAULT false,
    lead BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS locations(
   id SERIAL PRIMARY KEY,
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
   updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS ages(
    id SERIAL PRIMARY KEY,
    min INT,
    mix INT
);

CREATE TYPE status AS ENUM ('pending', 'open', 'not approved', 'did not occur', 'in review');
CREATE TABLE IF NOT EXISTS projects(
   id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   enabled BOOLEAN NOT NULL DEFAULT FALSE,
   status status NOT NULL,
   required TEXT NOT NULL,
   needed TEXT NOT NULL,
   leader_id TEXT NOT NULL,
   location_id BIGINT NOT NULL,
   start_time TIMESTAMPTZ NOT NULL,
   end_time TIMESTAMPTZ NOT NULL,
   category TEXT DEFAULT '',
   ages_id BIGINT,
   wheelchair BOOLEAN NOT NULL DEFAULT FALSE,
   short_description TEXT NOT NULL,
   long_description TEXT DEFAULT '',
   created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
   updated_at TIMESTAMPTZ,
   FOREIGN KEY (leader_id) REFERENCES accounts(id) ON DELETE RESTRICT,
   FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE RESTRICT,
   FOREIGN KEY (ages_id) REFERENCES ages(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS registrations(
    project_id BIGINT NOT NULL,
    account_id TEXT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE RESTRICT,
    UNIQUE (project_id, account_id)
);

CREATE TABLE IF NOT EXISTS project_skills(
     project_id BIGINT NOT NULL,
     skill TEXT NOT NULL,
     FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS project_tools(
     project_id BIGINT NOT NULL,
     tool TEXT NOT NULL,
     FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT
);

INSERT INTO accounts (id, first, last, email, cellphone, created_at, updated_at)
VALUES ('exampleid' , 'admin', 'user', 'adminserve@gmail.com', '303-947-7791',  LOCALTIMESTAMP, LOCALTIMESTAMP);

INSERT INTO locations (latitude, longitude, info, street, number, city, state, postal_code, formatted_address, created_at, updated_at)
VALUES ('39.5023509486161', '-104.87569755087917', 'Journey Church', 'Clydesdale Road', '9009', 'Castle Rock', 'Colorado', '80108', '9009 Clydesdale Rd, Castle Rock, CO 80108', LOCALTIMESTAMP, LOCALTIMESTAMP);

INSERT INTO projects (name, enabled, status, required, needed, start_time, end_time, leader_id, location_id, short_description, created_at, updated_at)
VALUES ('Base Project Example', true, 'open',  200, 159, LOCALTIMESTAMP, LOCALTIMESTAMP, 'exampleid', 1, 'Baseline Project Description', LOCALTIMESTAMP, LOCALTIMESTAMP);