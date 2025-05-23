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
                                        area TEXT,
                                        latitude DOUBLE PRECISION,
                                        longitude DOUBLE PRECISION,
                                        location_address TEXT,
                                        wheelchair_accessible BOOLEAN NOT NULL DEFAULT FALSE,
                                        serve_lead_id TEXT REFERENCES users(id),
                                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        status status NOT NULL DEFAULT 'open',
                                        website TEXT NOT NULL DEFAULT ''
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
                                             UNIQUE(user_id),
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


INSERT INTO projects (google_id, title, short_description, description, time, project_date,
                      max_capacity, area, latitude, longitude, serve_lead_id, wheelchair_accessible, location_address
) VALUES (1, 'Aging Resources Douglas County (Home 1)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (2, 'Aging Resources Douglas County (Home 2)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (3, 'Aging Resources Douglas County (Home 3)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (4, 'Aging Resources Douglas County (Home 4)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (5, 'Aging Resources Douglas County (Home 5)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (6, 'Aging Resources Douglas County (Home 6)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (7, 'Aging Resources Douglas County (Home 7)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (8, 'Aging Resources Douglas County (Home 8)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (9, 'Aging Resources Douglas County (Home 9)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (10, 'Aging Resources Douglas County (Home 10)', '', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', true, 'TBD. Communicated week prior to project.'), (11, 'Alternatives Pregnancy Center', '', 'Alternatives Pregnancy Center cares for Denver-area women and men in pregnancy-related crises and offers them a meaningful alternative to abortion.

Help us create a “baby shower in a bag” for every mom who chooses life by dropping off new, unopened baby items and gift cards to Journey in Castle Pines between 10:00 a.m. and 11:00 a.m.

Below is a list of items that are most needed:
·  $25 gift cards (to grocery stores like King Soopers and to retailers such as Walmart)
·  diapers (sizes N-6)
·  wipes
·  onesies (boys and girls, sizes 3–12 months)
·  sleepers (boys and girls, sizes 3–12 months)
·  baby socks
·  receiving blankets
·  bath towels and washcloths
·  baby lotion, soap,
·  baby books
·  small baby toys
·  pacifiers', '10 - 11am', '2025-07-12 00:00:00+00:00', 30, 'Englewood', 39.491482, -104.874878, 'example-user-123', true, '23 Inverness Way E. Suite 101B
Englewood , CO 80112'), (12, 'BackPack Society', '', 'The Backpack Society serves students who are struggling with food insecurity. Help us collect food donations to create easy-to-prepare weekend meals for students and their families. Donations can be dropped off at the Backpack Society at 213 W. County Line Road, Highlands Ranch, CO 80129, from 9:00 a.m. to 10:00 a.m.

Below is a list of food items that are most needed:
·  mac and cheese (single-serve cups and boxes)
·  canned beans
·  cereal
·  canned soup
·  canned pasta such as SpaghettiOs (large and small sizes)
·  ramen and Cup Noodles
·  pasta and pasta sauce
·  beef sticks or jerky (individual size)
·  tuna and chicken (cans and pouches)
·  bags of rice
·  fruit cups
·  applesauce pouches
·  granola bars, protein bars, and breakfast bars like Belvita
·  sports drinks and juices (small size)
·  snacks such as Chex Mix, chips, rice crisps, pretzels, Cheez-Its, and popcorn (individual size)
·  crackers of any variety (individual size and boxes for families)
·  jelly and jam (all flavors)
·  gluten-free items', '9 - 10am', '2025-07-12 00:00:00+00:00', 30, 'Highlands Ranch', 39.491482, -104.874878, 'example-user-123', true, '213 W. County Line Road
Highlands Ranch, CO 80129'), (13, 'Bin Blessed', '', 'Bin Blessed supports nonprofit There with Care in providing a wide range of meaningful services to families with children facing medical crises. This Serve Day, volunteers are needed to assemble movie cups containing candy and popcorn, and to decorate grocery bags that are used to deliver grocery items directly to families doorstep with a critically ill child.

This project will take place in Auditorium #2 at Journey Church in Castle Pines from 10:00 a.m. to 12:00 p.m. Reserved for families with small children so they can serve together!

Sign up to bring candy and popcorn using the Signup Genius link below:

https://www.signupgenius.com/go/5080548AFAD2AABF85-56498915-service', '10am - 12pm', '2025-07-12 00:00:00+00:00', 30, 'Journey Church (Castle Pines)', 39.491482, -104.874878, 'example-user-123', true, '9009 Clydesdale Road
Castle Pines, CO 80108'), (14, 'Box of Balloons', '', 'Box of Balloons is a nationwide nonprofit organization that celebrates children in need on their birthday. Join us to assemble and decorate birthday boxes that contain everything kids need to have a fantastic birthday party.

If you’d like to donate new birthday items, here is a generic list:
·  banners
·  tape
·  curling ribbon
·  streamers
·  gift cards
·  gift bags
·  party favors (for 6 kids)
·  tableware for 8 (plates, napkins, forks, cups)
·  party games for 6 kids
·  candles

We will assemble the birthday boxes at Journey Church in Castle Pines from 10:00 a.m. to 12:00 p.m. This is a great project for families with young children.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 20, 'Journey Church (Castle Pines)', 39.491482, -104.874878, 'example-user-123', true, '9009 Clydesdale Road
Castle Pines, CO 80108'), (15, 'Bridge of Hope', '', 'Bridge of Hope is a nonprofit committed to ending homelessness for single mothers and their children in the greater Denver area.
For Serve Day, we are collecting backpacks and school supplies for school-aged children in need.

New items needed for this project are:', '', '2025-07-12 00:00:00+00:00', 30, '', 39.491482, -104.874878, 'example-user-123', true, ''), (16, 'Colorado Feeding Kids (Shift 1)', '', 'Colorado Feeding Kids is a nonprofit whose mission is to provide nutritious food to impoverished people in Colorado and around the world. Serve Day volunteers will pack food using measuring cups, bags, and a heat press.

This project will take place at 14107 E. Exposition Ave., Aurora, CO 80012, from 11:00 a.m. to 1:00 p.m. Children ages 7 and up are welcome to participate alongside a parent.', '11am - 1pm ', '2025-07-12 00:00:00+00:00', 80, 'Aurora', 39.491482, -104.874878, 'example-user-123', true, '14107 E. Exposition Ave.
Aurora, CO 80012'), (17, 'Colorado Feeding Kids (Shift 2)', '', 'Colorado Feeding Kids is a nonprofit whose mission is to provide nutritious food to impoverished people in Colorado and around the world. Serve Day volunteers will pack food using measuring cups, bags, and a heat press.

This project will take place at 14107 E. Exposition Ave., Aurora, CO 80012, from 1:30 p.m. to 3:30 p.m. Children ages 7 and up are welcome to participate alongside a parent.', '1:30 - 3:30pm', '2025-07-12 00:00:00+00:00', 80, 'Aurora', 39.491482, -104.874878, 'example-user-123', true, '14107 E. Exposition Ave.
Aurora, CO 80012'), (18, 'Denver Rescue Mission (48th Ave Center)', '', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 7:15 a.m. to 9:00 a.m.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '7:15 - 9:00am', '2025-07-12 00:00:00+00:00', 6, 'Denver', 39.491482, -104.874878, 'example-user-123', true, '4600 E 48th Ave.
Denver, CO 80216'), (19, 'Denver Rescue Mission (48th Ave Center)', '', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the lunch meal service (including food preparation and distribution) from 11:15 a.m. to 1:00 p.m.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '11:15am - 1pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', true, '4600 E 48th Ave.
Denver, CO 80216'), (20, 'Denver Rescue Mission (48th Ave Center)', '', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the dinner meal service (including food preparation and distribution) from 4:45 p.m. to 6:30 p.m.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '4:45 - 6:30pm', '2025-07-12 00:00:00+00:00', 8, 'Denver', 39.491482, -104.874878, 'example-user-123', true, '4600 E 48th Ave.
Denver, CO 80216'), (21, 'Denver Rescue Mission (Lawrence St. Center)', '', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 6:45 a.m. to 8:45 a.m.

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '6:45 - 8:45am', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.491482, -104.874878, 'example-user-123', true, '2222 Lawrence St.
Denver, CO 80205'), (22, 'Denver Rescue Mission (Lawrence St. Center)', '', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing, following the breakfast service, from 7:30 AM to 9:00 AM. Volunteers must be 18 or older.

Located at 2222 Lawrence St, Denver, CO 80205.
', '7:30 - 9:00am', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', true, '2222 Lawrence St.
Denver, CO 80205'), (23, 'Denver Rescue Mission (Lawrence St. Center)', '', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing, following lunch service, from 12:30 PM to 2:00 PM. Volunteers must be 18 or older.

Located at 2222 Lawrence St, Denver, CO 80205.

', '12:30 - 2pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', true, '2222 Lawrence St.
Denver, CO 80205'), (24, 'Denver Rescue Mission (Lawrence St. Center)', '', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dinner meal prep from 2:30 PM to 3:30 PM. Volunteers must be 14 or older – 14 to 17 must be accompanied by an adult.

Located at 2222 Lawrence St, Denver, CO 80205.

', '2:30 - 3:30pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', true, '2222 Lawrence St.
Denver, CO 80205');