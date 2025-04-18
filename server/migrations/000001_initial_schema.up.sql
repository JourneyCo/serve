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
                                        time TEXT NOT NULL,
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

INSERT INTO projects (title, short_description, description, time, project_date,
                      max_capacity, location_name, latitude, longitude, lead_user_id, wheelchair_accessible, location_address
) VALUES (
             'Community Park Cleanup',
             'cleanup project',
             'Join us for a community park cleanup event! We will be cleaning up trash, planting flowers, and making general improvements to our local park. All supplies will be provided. Please wear comfortable clothes and bring water.',
             '9:00AM - 9:30AM',
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

CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          category TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS project_categories (
                                                  project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
                                                  category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
                                                  PRIMARY KEY (project_id, category_id)
);

CREATE TABLE IF NOT EXISTS ages(
                                   id SERIAL PRIMARY KEY,
                                   name TEXT NOT NULL
);

-- Insert age categories
INSERT INTO ages (name) VALUES
                            ('All ages'),
                            ('11-14'),
                            ('15-18'),
                            ('19-22'),
                            ('20s'),
                            ('30s'),
                            ('40s'),
                            ('50s'),
                            ('60s+');


CREATE TABLE IF NOT EXISTS project_ages (
                                            project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
                                            ages_id INTEGER REFERENCES ages(id) ON DELETE CASCADE,
                                            PRIMARY KEY (project_id, ages_id)
);

-- Insert categories
INSERT INTO categories (category) VALUES
                                      ('Arts & Crafts'),
                                      ('Community Outreach'),
                                      ('Español'),
                                      ('Family/Kid Friendly'),
                                      ('Food Prep & Distribution'),
                                      ('Indoors'),
                                      ('Landscaping'),
                                      ('Minor Home Repairs'),
                                      ('Outdoor'),
                                      ('Painting'),
                                      ('Prayer & Visitations'),
                                      ('Skilled Construction'),
                                      ('Sorting/Assembly'),
                                      ('Block Party'),
                                      ('Kids Ministry');

CREATE TABLE IF NOT EXISTS skills (
                                      id SERIAL PRIMARY KEY,
                                      name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS project_skills (
                                              project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
                                              skill_id INTEGER REFERENCES skills(id) ON DELETE CASCADE,
                                              PRIMARY KEY (project_id, skill_id)
);

-- Insert skills
INSERT INTO skills (name) VALUES
                              ('Carpentry'),
                              ('Communication'),
                              ('Construction'),
                              ('Cooking'),
                              ('Hospitality'),
                              ('Landscaping'),
                              ('Musical'),
                              ('Organizational'),
                              ('Painting'),
                              ('Photography');

CREATE TABLE IF NOT EXISTS supplies (
                                        id SERIAL PRIMARY KEY,
                                        name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS project_supplies (
                                                project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
                                                supply_id INTEGER REFERENCES supplies(id) ON DELETE CASCADE,
                                                PRIMARY KEY (project_id, supply_id)
);

-- Insert supplies
INSERT INTO supplies (name) VALUES
                                ('Bleach'),
                                ('Car Wash Supplies'),
                                ('Cleaning supplies'),
                                ('Craft supplies'),
                                ('Duct tape'),
                                ('Grilling Supplies'),
                                ('Landscape supplies'),
                                ('Nails'),
                                ('Paint supplies'),
                                ('Screws');

CREATE TABLE IF NOT EXISTS tools (
                                     id SERIAL PRIMARY KEY,
                                     name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS project_tools (
                                             project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
                                             tool_id INTEGER REFERENCES tools(id) ON DELETE CASCADE,
                                             PRIMARY KEY (project_id, tool_id)
);

-- Insert tools
INSERT INTO tools (name) VALUES
                             ('Chainsaw'),
                             ('Drill'),
                             ('Gloves'),
                             ('Hacksaw'),
                             ('Hammer'),
                             ('Hand saw'),
                             ('Hedge trimmers'),
                             ('Hoe'),
                             ('Ladder'),
                             ('Lawn mower'),
                             ('Leaf blower'),
                             ('Level'),
                             ('Miter saw'),
                             ('Nail gun'),
                             ('Paint gun'),
                             ('Pickaxe'),
                             ('Pitch fork'),
                             ('Pliers'),
                             ('Pressure washer'),
                             ('Pruners'),
                             ('Putty knife'),
                             ('Rake'),
                             ('Sander'),
                             ('Sawhorse'),
                             ('Screwdriver'),
                             ('Shovel'),
                             ('Sledgehammer'),
                             ('Socket set'),
                             ('Spanner wrench'),
                             ('Tape measure'),
                             ('Utility knife'),
                             ('Weedeater'),
                             ('Wheelbarrow'),
                             ('Wrench');
