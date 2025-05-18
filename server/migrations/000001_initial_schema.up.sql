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
           'servelead@example.com',
           'Doug',
           'DoesGood',
           '303-555-0123',
           true
       );

CREATE TYPE status AS ENUM ('pending', 'open', 'not_approved', 'did_not_occur', 'in_review');

CREATE TABLE IF NOT EXISTS projects (
                                        id SERIAL PRIMARY KEY,
                                        google_id INT,
                                        title TEXT NOT NULL,
                                        short_description TEXT NOT NULL,
                                        description TEXT NOT NULL,
                                        time TEXT NOT NULL,
                                        project_date timestamptz NOT NULL,
                                        max_capacity INTEGER NOT NULL,
                                        location_name TEXT,
                                        latitude DOUBLE PRECISION,
                                        longitude DOUBLE PRECISION,
                                        location_address TEXT,
                                        wheelchair_accessible BOOLEAN NOT NULL DEFAULT FALSE,
                                        serve_lead_id TEXT REFERENCES users(id),
                                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        status status NOT NULL DEFAULT 'open',
                                        website TEXT
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

DO $$
DECLARE
    serve_day TIMESTAMP;
BEGIN
    serve_day := NOW();

INSERT INTO projects (google_id, title, short_description, description, time, project_date,
                      max_capacity, location_name, latitude, longitude, serve_lead_id, wheelchair_accessible, location_address, website
) VALUES
      (1, 'Aging Resources Douglas County (Home 1)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (2, 'Aging Resources Douglas County (Home 2)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '1:30pm - 3:30pm', 8, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (3, 'Aging Resources Douglas County (Home 3)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (4, 'Aging Resources Douglas County (Home 4)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (5, 'Aging Resources Douglas County (Home 5)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (6, 'Aging Resources Douglas County (Home 6)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (7, 'Aging Resources Douglas County (Home 7)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (8, 'Aging Resources Douglas County (Home 8)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (9, 'Aging Resources Douglas County (Home 9)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (10, 'Aging Resources Douglas County (Home 10)', 'Elderly Support', 'Project scope to include yard work, landscaping, and gardening at the homes of senior adults. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. Location will be in Douglas County from 9:00 AM to 12:00 PM. Address will be provided closer to Serve Day. Any age welcome – kids must be accompanied by an adult. Volunteers younger than middle school will not be counted in the volunteer numbers.

**** All adults over the age of 18 will need to fill out an ARDC Volunteer Application.  Copy and paste this website: http://www.agingresourcesdougco.org/service-saturdays.html  ****

For more information on Aging Resources Douglas County, copy and paste this website:   https://www.agingresourcesdougco.org', 'www.example.com', '9:00 am - 12:00 pm', 7, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day', 39.491482, -104.874878, 'example-user-123', true, 'TBD (Douglas County)

Address of home to be provided closer to Serve Day'), (14, 'Box of Balloons', 'Underprivileged Children''s Birthday Parties', 'Box of Balloons is a nationwide non-profit organization that ensures underprivileged children are celebrated on their birthdays. This project will focus on putting together birthday boxes. Volunteers will decorate a banner and box and fill it with items to celebrate a child’s birthday.

Project will take place at Journey Church (Castle Pines location) from 10:00 AM to 12:00 PM. Project is best suited for families with young children.

*If you would like to donate, Box of Balloons utilizes banners, tape, curling ribbon, streamers, gifts, gift cards, party favors, tableware, party games, and candles.

                                                       www.boxofballoons.org', 'www.example.com', '10:00 am - 12:00 pm', 20, 'Journey Church
                                                       9009 Clydesdale Road
                                                       Castle Rock, CO 80108

', 39.491482, -104.874878, 'example-user-123', true, 'Journey Church
9009 Clydesdale Road
Castle Rock, CO 80108

'), (17, 'Denver Rescue Mission - 48th Ave Center ', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the dinner meal service (food prep and distribution) from 4:45 PM to 6:30 PM.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older – those aged 14 to 17 must be accompanied by an adult.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '4:45pm-6:30pm', 8, 'Denver

4600 E 48th Ave.
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

4600 E 48th Ave.
Denver, CO 80216'), (18, 'Denver Rescue Mission - 48th Ave Center ', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the lunch meal service (food prep and distribution) from 11:15 AM to 1:00 PM.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older – those aged 14 to 17 must be accompanied by an adult.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '11:15am-1:00pm', 2, 'Denver

4600 E 48th Ave.
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

4600 E 48th Ave.
Denver, CO 80216'), (19, 'Denver Rescue Mission - 48th Ave Center ', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the breakfast meal service (food prep and distribution) from 7:15 AM to 9:00 AM.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older – those aged 14 to 17 must be accompanied by an adult.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '7:15am-9:00am', 6, 'Denver

4600 E 48th Ave.
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

4600 E 48th Ave.
Denver, CO 80216'), (20, 'Denver Rescue Mission - Lawrence St Center', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing, following lunch service, from 12:30 PM to 2:00 PM. Volunteers must be 18 or older.

                                                       Located at 2222 Lawrence St, Denver, CO 80205.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '12:30-2:00pm', 2, 'Denver

2222 Lawrence St
Denver, CO 80205', 39.491482, -104.874878, 'example-user-123', true, 'Denver

2222 Lawrence St
Denver, CO 80205'), (21, 'Denver Rescue Mission - Lawrence St Center', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing, following the dinner service, from 5:30 PM to 7:00 PM. Volunteers must be 18 or older.

                                                       Located at 2222 Lawrence St, Denver, CO 80205.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '5:30pm-7:00pm', 2, 'Denver

2222 Lawrence St
Denver, CO 80205', 39.491482, -104.874878, 'example-user-123', true, 'Denver

2222 Lawrence St
Denver, CO 80205'), (22, 'Denver Rescue Mission - Lawrence St Center', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing, following the breakfast service, from 7:30 AM to 9:00 AM. Volunteers must be 18 or older.

                                                       Located at 2222 Lawrence St, Denver, CO 80205.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '7:30am-9:00am', 2, 'Denver

2222 Lawrence St
Denver, CO 80205', 39.491482, -104.874878, 'example-user-123', true, 'Denver

2222 Lawrence St
Denver, CO 80205'), (23, 'Denver Rescue Mission - Lawrence St Center', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with morning meal prep from 10:00 AM to 11:00 AM. Volunteers must be 14 or older – 14 to 17 must be accompanied by an adult.

Located at 2222 Lawrence St, Denver, CO 80205.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '10:00am-11:00am', 8, 'Denver

2222 Lawrence St
Denver, CO 80205', 39.491482, -104.874878, 'example-user-123', true, 'Denver

2222 Lawrence St
Denver, CO 80205'), (24, 'Denver Rescue Mission - Lawrence St Center', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dinner meal prep from 2:30 PM to 3:30 PM. Volunteers must be 14 or older – 14 to 17 must be accompanied by an adult.

Located at 2222 Lawrence St, Denver, CO 80205.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '2:30pm-3:30pm', 2, 'Denver

2222 Lawrence St
Denver, CO 80205', 39.491482, -104.874878, 'example-user-123', true, 'Denver

2222 Lawrence St
Denver, CO 80205'), (25, 'Denver Rescue Mission - Lawrence St Center', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the dinner meal service from 4:45 PM to 6:30 PM. Volunteers must be 14 or older – 14 to 17 must be accompanied by an adult.

Located at 2222 Lawrence St, Denver, CO 80205.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '4:45-6:30pm', 10, 'Denver

2222 Lawrence St
Denver, CO 80205', 39.491482, -104.874878, 'example-user-123', true, 'Denver

2222 Lawrence St
Denver, CO 80205'), (26, 'Denver Rescue Mission - Lawrence St Center', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the lunch meal service from 11:45 AM to 1:15 PM. Volunteers must be 14 or older – 14 to 17 must be accompanied by an adult.

Located at 2222 Lawrence St, Denver, CO 80205.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '11:45-1:15pm', 8, 'Denver

2222 Lawrence St
Denver, CO 80205', 39.491482, -104.874878, 'example-user-123', true, 'Denver

2222 Lawrence St
Denver, CO 80205'), (27, 'Denver Rescue Mission - The Crossing', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the breakfast meal service from 6:00 AM to 7:15 AM. Volunteers must be 12 or older – 12 to 17 must be accompanied by an adult.

Located at 6090 Smith Road, Denver, CO 80216.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '6:00am-7:15am', 3, 'Denver

6090 Smith Road
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

6090 Smith Road
Denver, CO 80216'), (28, 'Denver Rescue Mission- The Crossing', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve through participating in morning kitchen help from 8:00 AM to 10:00 AM. Volunteers must be 12 or older – 12 to 17 must be accompanied by an adult.

Located at 6090 Smith Road, Denver, CO 80216.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '8:00am-10:00am', 2, 'Denver

6090 Smith Road
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

6090 Smith Road
Denver, CO 80216'), (29, 'Denver Rescue Mission - The Crossing', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve through participating in afternoon kitchen help from 2:00 PM to 4:00 PM. Volunteers must be 12 or older – 12 to 17 must be accompanied by an adult.

Located at 6090 Smith Road, Denver, CO 80216.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '2:00pm-4:00pm', 3, 'Denver

6090 Smith Road
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

6090 Smith Road
Denver, CO 80216'), (30, 'Denver Rescue Mission - The Crossing', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the dinner meal service from 5:00 PM to 7:00 PM. Volunteers must be 12 or older – 12 to 17 must be accompanied by an adult.

Located at 6090 Smith Road, Denver, CO 80216.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '5:00pm-7:00pm', 6, 'Denver

6090 Smith Road
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

6090 Smith Road
Denver, CO 80216'), (31, 'Denver Rescue Mission - The Crossing', 'Serving Homeless/Those In Transition', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless/those in transition through participating in the lunch meal service from 11:00 AM to 1:00 PM. Volunteers must be 12 or older – 12 to 17 must be accompanied by an adult.

Located at 6090 Smith Road, Denver, CO 80216.

For more information on Denver Rescue Mission, copy and paste this website:   https://denverrescuemission.org', 'www.example.com', '11:00am-1:00pm', 6, 'Denver

6090 Smith Road
Denver, CO 80216', 39.491482, -104.874878, 'example-user-123', true, 'Denver

6090 Smith Road
Denver, CO 80216'), (38, 'SECOR 2-18-2025', 'Suburban Poverty Assistance', 'SECOR cares for people faced with Suburban Poverty.

Volunteers will assemble Food for Thought bags, pack dry goods boxes for the Mobile Market, repackage bulk items into smaller, single servings, and sort food donations. Volunteers must be 8 years old and must be accompanied by an adult. Closed-toe shoes are required. All volunteers must fill out a liability waiver here: https://forms.gle/ymctoRUt2iPfjSYQ8. Lastly, please stay home if you are feeling ill.

Address: 17151 Pine Ln., Parker, CO 80134. Please meet in the lobby, the building on the east end of the parking lot. Volunteering will take place from 10:00 to 11:30 a.m.

www.secorcares.com', 'www.example.com', '10:00 - 11:30 am', 8, 'Parker

17151 Pine Ln.
Parker, CO 80134                                             Meet in the lobby,the building on the east end of the parking lot.', 39.491482, -104.874878, 'example-user-123', true, 'Parker

17151 Pine Ln.
Parker, CO 80134                                             Meet in the lobby,the building on the east end of the parking lot.'), (42, 'We So They', 'Garage sale, organize, set up, tear down, help shoppers', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally; funds are raised through garage sales!

For this project, volunteers will help with a garage sale by organizing, setting up, tearing down, and assisting shoppers. Available in shifts from 7:30 AM to 10:30 AM, 10:30 AM to 1 PM, 1:00 PM to 4:00 PM, and 4:00 to 6:00 PM (need strong people for tear down). Volunteers can be 5+ with a parent, aside from the tear down timeslot.

Located in the Baptized Church parking lot at 880 Third Street, Castle Rock, CO. ', 'www.example.com', '7:30am-10:30am (10 people)    10:30-1:00 pm (10 people)   1:00-4:00pm (10 people)  4:00-6:00pm (10 strong people for tear down)', 40, '880 Third Street
Castle Rock, Co', 39.491482, -104.874878, 'example-user-123', true, '880 Third Street
Castle Rock, Co'), (45, 'Bridge of Hope  **2025 waiting to hear back from Veronika', 'Assemble goodie bags for their golf tournament fundraiser', 'Volunteers will assemble goodie bags for the Bridge of Hope golf fundraiser! Located at Journey Church Castle Pines Location from 10:00 AM to 12:00 PM. Any age welcome!

https://greaterdenver.bridgeofhopeinc.org/', 'www.example.com', '9-11', 10, 'Journey Church
9009 Clydesdale Road
Castle Rock, CO 80108

', 39.491482, -104.874878, 'example-user-123', true, 'Journey Church
9009 Clydesdale Road
Castle Rock, CO 80108

'), (46, 'Colorado Helping Hub', 'Local suburban poverty assistance', 'Calling all high schoolers ages 14-18! High school students will put on a food drive (1:00 PM to 4:00 PM) in Parker at multiple King Soopers locations. Students will pass out lists of items needed, pick up food items/supplies, and drop off at a central location where items will be packaged.

Please register on the Serve Day app but also register here: https://www.cohelpinghub.org/events

Signup Genius Link to sign up for task and location:
https://www.signupgenius.com/go/10C0B4EA8AF23A5FECF8-49521239-food

Addresses:
17031 Lincoln Ave, Parker, CO 80134,
17761 Cottonwood Dr, Parker, CO 80134,
12959 S Parker Rd, Parker, CO 80134
10901 S Parker Rd, Parker, CO 80134

https://www.cohelpinghub.org', 'www.example.com', '1:00pm-4:00pm', 30, '4 King Soopers locations:
17031 Lincoln Ave
Parker, CO 80134,

17761 Cottonwood Dr
Parker, CO 80134,

12959 S Parker Rd
Parker, CO 80134,

10901 S Parker Rd
Parker, CO 80134  ', 39.491482, -104.874878, 'example-user-123', true, '4 King Soopers locations:
17031 Lincoln Ave
Parker, CO 80134,

17761 Cottonwood Dr
Parker, CO 80134,

12959 S Parker Rd
Parker, CO 80134,

10901 S Parker Rd
Parker, CO 80134  ');

END $$;