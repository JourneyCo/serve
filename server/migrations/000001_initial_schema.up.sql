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
    max INT
);
INSERT INTO ages(min, max)  VALUES (18, 100);

CREATE TYPE status AS ENUM ('pending', 'open', 'not_approved', 'did_not_occur', 'in_review');
CREATE TABLE IF NOT EXISTS projects(
   id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   enabled BOOLEAN NOT NULL DEFAULT FALSE,
   status status NOT NULL,
   required TEXT NOT NULL,
   registered TEXT NOT NULL,
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

CREATE OR REPLACE FUNCTION update_total()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects
    SET registered = (SELECT COALESCE(SUM(quantity), 0)
                 FROM registrations
                 WHERE project_id = NEW.project_id)
    WHERE id = NEW.project_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trigger_update_total
    AFTER INSERT OR UPDATE OR DELETE ON registrations
    FOR EACH ROW
EXECUTE FUNCTION update_total();


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

INSERT INTO projects (name, enabled, status, required, registered, start_time, end_time, category, ages_id, leader_id, location_id, short_description, long_description, created_at, updated_at)
VALUES ('Base Project Example', true, 'open',  200, 0, LOCALTIMESTAMP, LOCALTIMESTAMP, 'open', 1, 'exampleid', 1, 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.', 'Nam pharetra neque et ornare consequat. Cras nec porttitor leo. Integer in pellentesque felis, sit amet mattis lorem. Nullam laoreet risus a diam posuere, id ultrices urna elementum. Vestibulum quis nisi sed lorem scelerisque dignissim ultrices in justo. Quisque pellentesque mattis justo eget ornare.', LOCALTIMESTAMP, LOCALTIMESTAMP);

INSERT INTO project_skills(project_id, skill) VALUES (1, 'carpentry');
INSERT INTO project_skills(project_id, skill) VALUES (1, 'painting');
INSERT INTO project_tools(project_id, tool) VALUES (1, 'hammer');
INSERT INTO project_tools(project_id, tool) VALUES (1, 'screwdriver');

