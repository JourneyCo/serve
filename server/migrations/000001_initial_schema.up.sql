CREATE TABLE IF NOT EXISTS users (
     id TEXT PRIMARY KEY,
     email TEXT NOT NULL DEFAULT '' UNIQUE,
     first_name TEXT NOT NULL DEFAULT '',
     last_name TEXT NOT NULL DEFAULT '',
     phone TEXT DEFAULT '',
     text_permission BOOLEAN DEFAULT FALSE,
     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO users (id, email, first_name, last_name, phone, text_permission)
VALUES (
   'example-user-123',
   'project.lead@example.com',
   'Doug',
   'DoesGood',
   '303-555-0123',
   true
);

CREATE TYPE status AS ENUM ('pending', 'open', 'not_approved', 'did_not_occur', 'in_review');

CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    short_description TEXT NOT NULL,
    description TEXT NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    project_date DATE NOT NULL,
    max_capacity INTEGER NOT NULL,
    location_name TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location_address TEXT,
    wheelchair_accessible BOOLEAN NOT NULL DEFAULT FALSE,
    lead_user_id TEXT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status status NOT NULL DEFAULT 'open'
);

INSERT INTO projects (title, short_description, description, start_time, end_time, project_date,
    max_capacity, location_name, latitude, longitude, lead_user_id, wheelchair_accessible, location_address
) VALUES (
             'Community Park Cleanup',
             'cleanup project',
             'Join us for a community park cleanup event! We will be cleaning up trash, planting flowers, and making general improvements to our local park. All supplies will be provided. Please wear comfortable clothes and bring water.',
             '09:00:00',
             '12:00:00',
             '2025-07-12',
             25,
             'Central Community Park',
             40.7128,
             -74.0060,
             'example-user-123',
             true,
             '123 Main Street, New York, NY 10001'
         );

CREATE TABLE IF NOT EXISTS registrations (
     id SERIAL PRIMARY KEY,
     user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
     project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
     status TEXT NOT NULL DEFAULT 'registered',
     guest_count INTEGER NOT NULL DEFAULT 0,
     lead_interest BOOLEAN NOT NULL DEFAULT FALSE,
     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
     UNIQUE(user_id, project_id),
     FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT,
     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS needs (
     id SERIAL PRIMARY KEY,
     name TEXT NOT NULL UNIQUE,
        type TEXT NOT NULL DEFAULT 'tool',
     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_needs (
     project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
     need_id INTEGER REFERENCES needs(id) ON DELETE CASCADE,
     PRIMARY KEY (project_id, need_id)
);

CREATE TABLE IF NOT EXISTS ages(
    id SERIAL PRIMARY KEY,
    min INT,
    max INT
);
INSERT INTO ages(min, max)  VALUES (18, 100);


CREATE TABLE IF NOT EXISTS project_ages (
     project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
     ages_id INTEGER REFERENCES ages(id) ON DELETE CASCADE,
     PRIMARY KEY (project_id, ages_id)
);



INSERT INTO needs(id, name) VALUES (1, 'carpentry');
INSERT INTO needs(id, name) VALUES (2, 'painting');
INSERT INTO project_needs(project_id, need_id) VALUES (1, 1);
INSERT INTO project_needs(project_id, need_id) VALUES (1, 2);


