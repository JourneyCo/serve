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
                                        description TEXT NOT NULL,
                                        time TEXT NOT NULL,
                                        project_date timestamptz NOT NULL,
                                        max_capacity INTEGER NOT NULL,
                                        area TEXT,
                                        latitude DOUBLE PRECISION,
                                        longitude DOUBLE PRECISION,
                                        location_address TEXT,
                                        serve_lead_id TEXT REFERENCES users(id),
                                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        status status NOT NULL DEFAULT 'open',
                                        website TEXT NOT NULL DEFAULT '',
                                        ages TEXT NOT NULL DEFAULT 'All Ages'
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

CREATE TABLE IF NOT EXISTS types (
                                          id INTEGER PRIMARY KEY,
                                          type TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS project_types (
                                                  project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
                                                  type_id INTEGER REFERENCES types(id) ON DELETE CASCADE,
                                                  PRIMARY KEY (project_id, type_id)
);

INSERT INTO projects (google_id, title, description, time, project_date, max_capacity, area, latitude, longitude, serve_lead_id, location_address, website, ages) VALUES (1, 'Aging Resources Douglas County (Home 1)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (2, 'Aging Resources Douglas County (Home 2)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (3, 'Aging Resources Douglas County (Home 3)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (4, 'Aging Resources Douglas County (Home 4)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (5, 'Aging Resources Douglas County (Home 5)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (6, 'Aging Resources Douglas County (Home 6)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (7, 'Aging Resources Douglas County (Home 7)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (8, 'Aging Resources Douglas County (Home 8)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (9, 'Aging Resources Douglas County (Home 9)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (10, 'Aging Resources Douglas County (Home 10)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc.

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers.

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (11, 'Alternatives Pregnancy Center', 'Alternatives Pregnancy Center cares for Denver-area women and men in pregnancy-related crises and offers them a meaningful alternative to abortion.

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
·  pacifiers', '10 - 11am', '2025-07-12 00:00:00+00:00', 30, 'Englewood', 39.566699, -104.859965, 'example-user-123', '23 Inverness Way E. Suite 101B
Englewood , CO 80112', 'www.youhavealternatives.org', 'All Ages'), (12, 'BackPack Society', 'The Backpack Society serves students who are struggling with food insecurity. Help us collect food donations to create easy-to-prepare weekend meals for students and their families. Donations can be dropped off at the Backpack Society at 213 W. County Line Road, Highlands Ranch, CO 80129, from 9:00 a.m. to 10:00 a.m.

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
·  gluten-free items', '9 - 10am', '2025-07-12 00:00:00+00:00', 30, 'Highlands Ranch', 39.567481, -104.991901, 'example-user-123', '213 W. County Line Road
Highlands Ranch, CO 80129', 'www.backpacksociety.org', 'All Ages'), (13, 'Bin Blessed', 'Bin Blessed supports nonprofit There with Care in providing a wide range of meaningful services to families with children facing medical crises. This Serve Day, volunteers are needed to assemble movie cups containing candy and popcorn, and to decorate grocery bags that are used to deliver grocery items directly to families doorstep with a critically ill child.

This project will take place in Auditorium #2 at Journey Church in Castle Pines from 10:00 a.m. to 12:00 p.m. Reserved for families with small children so they can serve together!

Sign up to bring candy and popcorn using the Signup Genius link below:

https://www.signupgenius.com/go/5080548AFAD2AABF85-56498915-service', '10am - 12pm', '2025-07-12 00:00:00+00:00', 30, 'Journey Church (Castle Pines)', 39.491739, -104.874617, 'example-user-123', '9009 Clydesdale Road
Castle Pines, CO 80108', 'www.binblessed.com', 'Families with Small Children ONLY'), (14, 'Box of Balloons', 'Box of Balloons is a nationwide nonprofit organization that celebrates children in need on their birthday. Join us to assemble and decorate birthday boxes that contain everything kids need to have a fantastic birthday party.

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

We will assemble the birthday boxes at Journey Church in Castle Pines from 10:00 a.m. to 12:00 p.m. This is a great project for families with young children.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 20, 'Journey Church (Castle Pines)', 39.491739, -104.874617, 'example-user-123', '9009 Clydesdale Road
Castle Pines, CO 80108', 'www.boxofballoons.org', 'Families with Young Children'), (15, 'Bridge of Hope', 'Bridge of Hope is a nonprofit committed to ending homelessness for single mothers and their children in the greater Denver area.
For Serve Day, we are collecting backpacks and school supplies for school-aged children in need.

New items needed for this project are:', '', '2025-07-12 00:00:00+00:00', 30, '', 39.491482, -104.874878, 'example-user-123', '', 'www.greaterdenver.bridgeofhopeinc.org', 'All Ages'), (16, 'Colorado Feeding Kids (Shift 1)', 'Colorado Feeding Kids is a nonprofit whose mission is to provide nutritious food to impoverished people in Colorado and around the world. Serve Day volunteers will pack food using measuring cups, bags, and a heat press.

This project will take place at 14107 E. Exposition Ave., Aurora, CO 80012, from 11:00 a.m. to 1:00 p.m. Children ages 7 and up are welcome to participate alongside a parent.', '11am - 1pm ', '2025-07-12 00:00:00+00:00', 80, 'Aurora', 39.704600, -104.824789, 'example-user-123', '14107 E. Exposition Ave.
Aurora, CO 80012', 'www.cofeedingkids.org', '7 Years and Older'), (17, 'Colorado Feeding Kids (Shift 2)', 'Colorado Feeding Kids is a nonprofit whose mission is to provide nutritious food to impoverished people in Colorado and around the world. Serve Day volunteers will pack food using measuring cups, bags, and a heat press.

This project will take place at 14107 E. Exposition Ave., Aurora, CO 80012, from 1:30 p.m. to 3:30 p.m. Children ages 7 and up are welcome to participate alongside a parent.', '1:30 - 3:30pm', '2025-07-12 00:00:00+00:00', 80, 'Aurora', 39.704600, -104.824789, 'example-user-123', '14107 E. Exposition Ave.
Aurora, CO 80012', 'www.cofeedingkids.org', '7 Years and Older'), (18, 'Denver Rescue Mission (48th Ave Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 7:15 a.m. to 9:00 a.m.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '7:15 - 9:00am', '2025-07-12 00:00:00+00:00', 6, 'Denver', 39.782790, -104.934248, 'example-user-123', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '14 Years and Older'), (19, 'Denver Rescue Mission (48th Ave Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the lunch meal service (including food preparation and distribution) from 11:15 a.m. to 1:00 p.m.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '11:15am - 1pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.782790, -104.934248, 'example-user-123', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '14 Years and Older'), (20, 'Denver Rescue Mission (48th Ave Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the dinner meal service (including food preparation and distribution) from 4:45 p.m. to 6:30 p.m.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '4:45 - 6:30pm', '2025-07-12 00:00:00+00:00', 8, 'Denver', 39.782790, -104.934248, 'example-user-123', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '14 Years and Older'), (21, 'Denver Rescue Mission (Lawrence St. Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 6:45 a.m. to 8:45 a.m.

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '6:45 - 8:45am', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.754798, -104.988222, 'example-user-123', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '14 Years and Older'), (22, 'Denver Rescue Mission (Lawrence St. Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing following the breakfast service from 7:30 a.m. to 9:00 a.m.

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 18 or older.', '7:30 - 9:00am', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.754798, -104.988222, 'example-user-123', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '18 Years and Older'), (23, 'Denver Rescue Mission (Lawrence St. Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing following the lunch service from 12:30 p.m. to 2:00 p.m.

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 18 or older.', '12:30 - 2pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.754798, -104.988222, 'example-user-123', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '18 Years and Older'), (24, 'Denver Rescue Mission (Lawrence St. Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help prepare the dinner service from 2:30 p.m. to 3:30 p.m.

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '2:30 - 3:30pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.754798, -104.988222, 'example-user-123', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '14 Years and Older'), (25, 'Denver Rescue Mission (Lawrence St. Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the dinner meal service (including food preparation and distribution) from 4:45 p.m. to 6:30 p.m.

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '4:45 - 6:30pm', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.754798, -104.988222, 'example-user-123', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '14 Years and Older'), (26, 'Denver Rescue Mission (Lawrence St. Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing following the dinner service from 5:30 p.m. to 7:00 p.m.

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 18 or older.', '5:30 - 7:00pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.754798, -104.988222, 'example-user-123', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '18 Years and Older'), (27, 'Denver Rescue Mission (48th Ave Center)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will help with dishwashing following the lunch service from 12:00 p.m. to 1:30 p.m.

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 18 or older.', '12 - 1:30pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.782790, -104.934248, 'example-user-123', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '18 Years and Older'), (28, 'Denver Rescue Mission (The Crossing)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 6:00 a.m. to 7:15 a.m.

Denver Rescue Mission is located at 6090 Smith Rd., Denver, CO 80216. Volunteers must be 12 or older; those aged 12 to 17 must be accompanied by an adult.', '6 - 7:15am', '2025-07-12 00:00:00+00:00', 3, 'Denver', 39.773293, -104.918409, 'example-user-123', '6090 Smith Road
Denver, CO 80216', 'www.denverrescuemission.org', '12 Years and Older'), (29, 'Denver Rescue Mission (The Crossing)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will assist with afternoon kitchen duties from 2:00 p.m. to 4:00 p.m.

Denver Rescue Mission is located at 6090 Smith Rd., Denver, CO 80216. Volunteers must be 12 or older; those aged 12 to 17 must be accompanied by an adult.', '2 - 4pm', '2025-07-12 00:00:00+00:00', 3, 'Denver', 39.773293, -104.918409, 'example-user-123', '6090 Smith Road
Denver, CO 80216', 'www.denverrescuemission.org', '12 Years and Older'), (30, 'Denver Rescue Mission (The Crossing)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives.

Volunteers will serve the homeless and those in transition through participating in the dinner meal service (including food preparation and distribution) from 5:00 p.m. to 7:00 p.m.

Denver Rescue Mission is located at 6090 Smith Rd., Denver, CO 80216. Volunteers must be 12 or older; those aged 12 to 17 must be accompanied by an adult.', '5 - 7pm', '2025-07-12 00:00:00+00:00', 6, 'Denver', 39.773293, -104.918409, 'example-user-123', '6090 Smith Road
Denver, CO 80216', 'www.denverrescuemission.org', '12 Years and Older'), (31, 'Denver Street School', 'The Denver Street School is a high school program that brings the love of Christ to students who are struggling academically. It provides a quality education as well as support to help kids thrive.

Volunteers will assist with cleaning and organizing the teacher resource room from 9:00 a.m. to 12:00 p.m.

The Denver Street School is located at 1380 Ammons St., Lakewood, CO 80214. Children are welcome to participate alongside a parent.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 5, 'Lakewood', 39.735577, -105.087019, 'example-user-123', '1380 Ammons St
Lakewood, CO 80214', 'www.denverstreetschool.org', 'All Ages'), (32, 'Denver Street School', 'The Denver Street School is a high school program that brings the love of Christ to students who are struggling academically. It provides a quality education as well as support to help kids thrive.

Volunteers will assist with painting the teacher lounge from 9:00 a.m. to 12:00 p.m.

The Denver Street School is located at 1380 Ammons St., Lakewood, CO 80214. Children are welcome to participate alongside a parent.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 5, 'Lakewood', 39.735577, -105.087019, 'example-user-123', '1380 Ammons St
Lakewood, CO 80214', 'www.denverstreetschool.org', '18 Years and Older'), (33, 'Food Bank of the Rockies', 'Food Bank of the Rockies provides food and other necessities to people facing hunger in Colorado and Wyoming.

Volunteers will report to the main warehouse located at 10700 E 47th Ave., Denver, CO 80239, where they will inspect, clean, sort, and repack perishable and nonperishable food items for distribution to Hunger Relief Partners from 8:45 a.m. to 12:00 p.m. All group participants must be signed up two weeks prior to Serve Day, by June 28, 2025. Children ages 10 and up are welcome to participate alongside an adult chaperone (see chaperone requirements below).

Food Bank of the Rockies requires each volunteer to register and sign a waiver ahead of volunteering and will allow them to track their volunteer hours with Food Bank of the Rockies.

To register, visit https://vhub.at/1S4HE8D. Click on the date and the green “Sign up” button next to the event on the page. Then you will be prompted to either log in or create an account. You do not need a join code and can skip this step. Once you have created an account and registered for the event, you will receive a confirmation email.

If your group includes 10- and 11-year-olds, please have 1 chaperone for every 3 youths.
If your group has 12- to 14-year-olds, please have 1 chaperone for every 4 youths.
If your group has 15- to 17-year-olds, please have 1 chaperone for every 5 youths.
Peers 18 years of age do not count as chaperones.

To better prepare for volunteering, please watch this short Welcome Video and Safety Video.

If you have any questions, please reach out to volunteer@foodbankrockies.org.', '8:45am - 12pm', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.777716, -104.863715, 'example-user-123', '10700 E. 45th Ave.
Denver, CO 80239', 'www.foodbankrockies.org', '10 Years and Older'), (34, 'Food Bank of the Rockies', 'Food Bank of the Rockies provides food and other necessities to people facing hunger in Colorado and Wyoming.

Volunteers will report to the main warehouse located at 10700 E 47th Ave., Denver, CO 80239, where they will inspect, clean, sort, and repack perishable and nonperishable food items for distribution to Hunger Relief Partners from 1:00 p.m. to 4:00 p.m. All group participants must be signed up two weeks prior to Serve Day, by June 28, 2025. Children ages 10 and up are welcome to participate alongside an adult chaperone (see chaperone requirements below).

Food Bank of the Rockies requires each volunteer to register and sign a waiver ahead of volunteering and will allow them to track their volunteer hours with Food Bank of the Rockies.

To register, visit https://vhub.at/1S4HE8D. Click on the date and the green “Sign up” button next to the event on the page. Then you will be prompted to either log in or create an account. You do not need a join code and can skip this step. Once you have created an account and registered for the event, you will receive a confirmation email.

If your group includes 10- and 11-year-olds, please have 1 chaperone for every 3 youths.
If your group has 12- to 14-year-olds, please have 1 chaperone for every 4 youths.
If your group has 15- to 17-year-olds, please have 1 chaperone for every 5 youths.
Peers 18 years of age do not count as chaperones.

To better prepare for volunteering, please watch this short Welcome Video and Safety Video.

If you have any questions, please reach out to volunteer@foodbankrockies.org.', '1 - 4pm', '2025-07-12 00:00:00+00:00', 24, 'Denver', 39.777716, -104.863715, 'example-user-123', '10700 E. 45th Ave.
Denver, CO 80239', 'www.foodbankrockies.org', '10 Years and Older'), (35, 'Help and Hope Center', 'Help & Hope Center supports Douglas and Elbert County residents in financial distress, allowing them to navigate troublesome times with dignity.

Volunteers will pack Strive to Thrive bags, sort donations in receiving area, perform light cleaning projects, and assist with outdoor work such as weeding and general beautification. Volunteering will take place from 9:00 a.m. to 11:00 a.m. at 1638 Park St., Castle Rock, CO 80109.

Be advised that work may be performed outside and/or dirty, so participants should wear weather-appropriate wear clothing and closed-toe shoes. Volunteers must be at least 13 years old and accompanied by an adult.', '9 - 11am', '2025-07-12 00:00:00+00:00', 20, 'Castle Rock', 39.385684, -104.867875, 'example-user-123', '1638 Park St.
Castle Rock, CO 80109', 'www.helpandhopecenter.org', '13 Years and Older'), (36, 'Hope''s Promise (Project #1)', 'Hope’s Promise strives to strengthen families through ethical, Christ-centered foster care, adoption, and global orphan care.

Volunteers will make welcoming meals for foster children and their new families. Drop off welcoming meals at Hope’s Promise between 10:00 a.m. and 12:00 p.m. at 1585 S Perry St., #E, Castle Rock, CO 80104.

Make a freezer meal for a family of 4-5 people. Label meal with reheating instructions and a list ingredients. Bring frozen cookie dough for dessert, too. Any volunteers who would like a tour of Hope’s Promise can arrive early 9:45 a.m.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 20, 'Castle Rock', 39.357237, -104.861606, 'example-user-123', '1585 S Perry St # E
Castle Rock, CO 80104', 'www.hopespromise.com', 'All Ages'), (37, 'Hope''s Promise (Project #2)', 'Hope’s Promise strives to strengthen families through ethical, Christ-centered foster care, adoption, and global orphan care.

Volunteers will help prepare for an upcoming golf tournament by helping to assemble silent auction items and baskets along with creating signage and banquet decor.

Volunteering will take place at 1585 S Perry St., #E, Castle Rock, CO 80104 from 10:00 a.m. to 12:00 p.m.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Castle Rock', 39.357237, -104.861606, 'example-user-123', '1585 S Perry St # E
Castle Rock, CO 80104', 'www.hopespromise.com', '10 Years and Older'), (38, 'Hope''s Promise (Project #3)', 'Hope’s Promise strives to strengthen families through ethical, Christ-centered foster care, adoption, and global orphan care.

Volunteers will write encouraging, uplifting notes to parents, caregivers, and foster parents both internationally and in Colorado. Please bring a package of stationary or pretty note cards for this project.

Volunteering will take place at 1585 S Perry St., #E, Castle Rock, CO 80104 from 10:00 a.m. to 12:00 p.m.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 15, 'Castle Rock', 39.357237, -104.861606, 'example-user-123', '1585 S Perry St # E
Castle Rock, CO 80104', 'www.hopespromise.com', 'All Ages'), (39, 'His Love Fellowship Church', 'His Love Fellowship is a Bible-based church in west Denver.

Volunteers will help stain railroad tiles around the trees and paint over graffiti. His Love Fellowship is located at 910 Kalamath St., Denver, CO 80204.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.730739, -104.999721, 'example-user-123', '910 Kalamath St.
Denver, CO 80204', 'www.hislovefellowship.com', '12 Years and Older'), (40, 'SECOR', 'SECOR cares for people faced with suburban poverty by offering a Free Food Market where guests can shop for food and other necessities. Most guests are able to go home with enough food for 10-14 days.


This project will involve both market and warehouse service at SECOR, located at 17151 Pine Ln., Parker, CO 80134, from 8:45 a.m. to 11:00 a.m.

Market volunteers will help guests get started after they check in with the front desk, man various stations in the market as guests shop, and manage the cart return. Warehouse volunteers will help sort items that come in on the truck, restock market shelves, and help clean the market and warehouse toward the end of the shift. They may also assist with packing boxes, sorting donations, and assembling items. The first shift of volunteers will show the second shift the ropes, as SECOR will not have enough staff to handle that midday while guests are shopping.


Volunteers must be 8 years old and accompanied by an adult (groups of minors need to have a 1 adult to 6 minors ratio). Closed-toe shoes are required. All volunteers must fill out a liability waiver: https://forms.gle/ymctoRUt2iPfjSYQ8.
Please meet in the lobby of the building on the east end of the parking lot. Bring a lock if you want to lock up any valuables. Lastly, please stay home if you are feeling ill.', '8:45 - 11am', '2025-07-12 00:00:00+00:00', 10, 'Parker', 39.543951, -104.790072, 'example-user-123', '17151 Pine Lane
Parker, CO 80134', 'www.secorcares.com', '8 Years and Older'), (41, 'SECOR', 'SECOR cares for people faced with suburban poverty by offering a Free Food Market where guests can shop for food and other necessities. Most guests are able to go home with enough food for 10-14 days.


This project will involve both market and warehouse service at SECOR, located at 17151 Pine Ln., Parker, CO 80134, from 10:45 a.m. to 1:15 p.m.

Market volunteers will help guests get started after they check in with the front desk, man various stations in the market as guests shop, and manage the cart return. Warehouse volunteers will help sort items that come in on the truck, restock market shelves, and help clean the market and warehouse toward the end of the shift. They may also assist with packing boxes, sorting donations, and assembling items. The first shift of volunteers will show the second shift the ropes, as SECOR will not have enough staff to handle that midday while guests are shopping.


Volunteers must be 8 years old and accompanied by an adult (groups of minors need to have a 1 adult to 6 minors ratio). Closed-toe shoes are required. All volunteers must fill out a liability waiver: https://forms.gle/ymctoRUt2iPfjSYQ8.
Please meet in the lobby of the building on the east end of the parking lot. Bring a lock if you want to lock up any valuables. Lastly, please stay home if you are feeling ill.', '10:45am - 1:15pm', '2025-07-12 00:00:00+00:00', 10, 'Parker', 39.543951, -104.790072, 'example-user-123', '17151 Pine Lane
Parker, CO 80134', 'www.secorcares.com', '8 Years and Older'), (42, 'Sweet Dream In A Bag', 'Assembling bags with bedding that will be distributed to children in poverty locally and internationally.

"Sweet Dream in a Bag has provided cozy bedding to children in need since 2010. For this project, volunteers will assemble bags with bedding to be distributed to children in need, both locally and internationally. All ages are welcome. Event will run from 9:00 AM until 11:00 AM at 5933 S. Fairfield Street, Littleton, CO 80120.


"All ages (Kids accompanied by adult).  Children under the age of 5 will be more of a ""helper"" to their parents."', '9 - 11am', '2025-07-12 00:00:00+00:00', 20, 'Littleton', 39.609193, -104.990684, 'example-user-123', '5933 S. Fairfield Street
Littleton, CO 80120', 'www.sweetdreaminabag.org', 'All Ages'), (43, 'Vitalant Blood Bank', 'Blood Donations

This project will need 56 people willing to donate blood between 9 am and 1 pm. (Will Send link and QR code to Dave and Cory for website sign up.) The links must be included on the website sign up.

Signup Link:  https://donors.vitalant.org/dwp/portal/dwa/appointment/guest/phl/timeSlotsExtr?token=YejHtw9ohlDWKytYWH8lXTMxzDNXba93MJ5XkjWk610%3D

Eligibility Requirements: https://www.vitalant.org/eligibility

"Volunteers (16+ and 110+ pounds) are able to donate blood between 9 am and 1 pm at the Journey Church Castle Pines location in the upstairs atrium. Please see important links below.

Signup Link:https://donors.vitalant.org/dwp/portal/dwa/appointment/guest/phl/timeSlotsExtr?token=YejHtw9ohlDWKytYWH8lXWexFLCpLg9%2Bu7xOmseb9Qs%3D

Eligibility Requirements: https://www.vitalant.org/eligibility

https://www.vitalant.org/"', '9am - 1pm', '2025-07-12 00:00:00+00:00', 56, 'Journey Church (Castle Pines)', 39.491739, -104.874617, 'example-user-123', '9009 Clydesdale Road
Castle Pines, CO 80108', 'www.vitalant.org', '16 Years and Older'), (44, 'WeSoThey (Shift 1)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally; funds are raised through garage sales!

For this project, volunteers will help with a garage sale by organizing, setting up, tearing down, and assisting shoppers. Available in shifts from 7:30 AM to 10:30 AM, 10:30 AM to 1 PM, 1:00 PM to 4:00 PM, and 4:00 to 6:00 PM (need strong people for tear down). Volunteers can be 5+ with a parent, aside from the tear down timeslot.

3 hour shifts, multiple shifts, Family friendly ages 5+.  Set up items, tables, organize, assist shoppers, tear down, lift heavy items.', '7:30 - 10:30am', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.372236, -104.853072, 'example-user-123', '880 Third Street
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (45, 'WeSoThey (Shift 2)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally; funds are raised through garage sales!

For this project, volunteers will help with a garage sale by organizing, setting up, tearing down, and assisting shoppers. Available in shifts from 7:30 AM to 10:30 AM, 10:30 AM to 1 PM, 1:00 PM to 4:00 PM, and 4:00 to 6:00 PM (need strong people for tear down). Volunteers can be 5+ with a parent, aside from the tear down timeslot.

3 hour shifts, multiple shifts, Family friendly ages 5+.  Set up items, tables, organize, assist shoppers, tear down, lift heavy items.', '10:30am - 1pm', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.372236, -104.853072, 'example-user-123', '880 Third Street
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (46, 'WeSoThey (Shift 3)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally; funds are raised through garage sales!

For this project, volunteers will help with a garage sale by organizing, setting up, tearing down, and assisting shoppers. Available in shifts from 7:30 AM to 10:30 AM, 10:30 AM to 1 PM, 1:00 PM to 4:00 PM, and 4:00 to 6:00 PM (need strong people for tear down). Volunteers can be 5+ with a parent, aside from the tear down timeslot.

3 hour shifts, multiple shifts, Family friendly ages 5+.  Set up items, tables, organize, assist shoppers, tear down, lift heavy items.', '1 - 4pm', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.372236, -104.853072, 'example-user-123', '880 Third Street
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (47, 'WeSoThey (Shift 4)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally; funds are raised through garage sales!

For this project, volunteers will help with a garage sale by organizing, setting up, tearing down, and assisting shoppers. Available in shifts from 7:30 AM to 10:30 AM, 10:30 AM to 1 PM, 1:00 PM to 4:00 PM, and 4:00 to 6:00 PM (need strong people for tear down). Volunteers can be 5+ with a parent, aside from the tear down timeslot.

3 hour shifts, multiple shifts, Family friendly ages 5+.  Set up items, tables, organize, assist shoppers, tear down, lift heavy items.', '4 - 6pm', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.372236, -104.853072, 'example-user-123', '880 Third Street
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (48, 'Leman Academy (Journey Parker)', 'Volunteers will have the opportunity to love on the school that hosts our Parker location by cleaning up weeds and trash around the school, cleaning up the curb areas (sweeping rocks from curbs/trash from curbs), sweeping wood chips from the southside sidewalk back into the play areas, and disposing of any broken/unusable cones. The times will be 9 am to 12 pm. All ages are welcome, and any kids must be accompanied by an adult.

"Weeds/trash in the rocks around the school and in the medians. Curb clean up (sweeping rocks from curbs/trash from curbs). Sweeping the wood chips from the southside sidewalk back into the play areas
"', '9am - 12pm', '2025-07-12 00:00:00+00:00', 20, 'Parker', 39.477281, -104.762508, 'example-user-123', '19560 Stroh Rd.
Parker, CO 80134', 'N/A', 'All Ages'), (49, 'Barksdale Family Home Project', 'Physical Labor- Cutting /brush/ scrub oak and putting them in piles.

Ages that are able to do physical labor/yardwork.No small children please, due to the use of a lot of chainsaws.

Helping Journey Church members John and ? Barksdale with cutting down scrub oak and putting them into piles.

The Serve Day Lead will email volunteers his address as the date gets closer.

Volunteers will assist with much-needed yard work, cutting brush/scrub oak.  Please wear long sleeves, a hat;bring sunscreen, and gloves. Cool drinks will be provided. The project will take place from 8:00 AM to 11:00 PM

"Wear long sleeves, bring gloves, hat and sunscreen. John will provide cool drinks for all volunteers.Please bring cutting tools-chainsaws, and any tool that can cut brush.
 "', '8 - 11am', '2025-07-12 00:00:00+00:00', 20, 'Castle Rock', 39.407939, -104.856225, 'example-user-123', 'Address will be provided separately for those registered', 'N/A', '8 Years and Older'), (50, 'Open Door Ministries', 'ODM Cornerstone home is for disabled men with just 4 bedrooms.

Deep cleaning of main living areas of the home-living room, kitchen, hallways.  Light landscape work, planting some pots in front of the home, weeding the mulch beds.

Wear long sleeves, bring gloves, hat, sunscreen and water, cleaning supplies/rags.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 15, 'Denver', 39.756680, -104.978550, 'example-user-123', '2754 Stout St.
Denver, CO 80205', 'www.odmdenver.org', 'All Ages'), (51, 'Parker Task Force', 'Volunteers will support the Parker Task Force by first touring the Food Bank facility and recieving training.  Then, they will assist market clients when they come in to shop from 9-1.  ', '8:30am - 1pm', '2025-07-12 00:00:00+00:00', 4, 'Parker', 39.525121, -104.767662, 'example-user-123', '19105 Longs Way
Parker, CO 80134', 'www.parkertaskforce.org', '18 Years and Older'), (52, 'Affinity Ranch (Project #1)', 'Affinity Ranch serves the community, including veterans and those with differing abilities, through healing with animals. They currenlty have 18 horses and they take great care to make sure they are healthy and happy.  Volunteers will work with ranch staff to build and donate a new hay feeding trough for the horses.

Build a hay feeding trough for the horses.  This will be done in conjunction with ranch staff.

"Ranch provides supplies, bring gloves, safety glasses, wear closed toe work shoes/boots.  Long sleeves and a hat are suggested.  Waiver must be signed ahead of time.
https://prayinghandsranch.org/ranchapps/PublicForms/VolunteerGroupRegister.php?GrpName=JOURNEYCHURCH

"', '9am - 1pm', '2025-07-12 00:00:00+00:00', 5, 'Parker', 39.433931, -104.667828, 'example-user-123', '11892 Hilltop Road
Parker, 80134', 'www.affinityranch.org', '18 Years and Older');
INSERT INTO types (id, type)  VALUES (5, 'Community Outreach'), (7, 'Painting'), (1, 'Yard Work'), (3, 'Collection/Drop Off'), (4, 'Sorting/Assembly'), (6, 'Food Prep/Distribution'), (8, 'Minor Indoor Repairs'), (9, 'Arts & Crafts'), (10, 'Construction'), (2, 'Kid Friendly');
INSERT INTO project_types (project_id, type_id)  VALUES (1, '1'), (1, '2'), (2, '1'), (2, '2'), (3, '1'), (3, '2'), (4, '1'), (4, '2'), (5, '1'), (5, '2'), (6, '1'), (6, '2'), (7, '1'), (7, '2'), (8, '1'), (8, '2'), (9, '1'), (9, '2'), (10, '1'), (10, '2'), (11, '3'), (11, '2'), (12, '3'), (12, '2'), (13, '4'), (13, '5'), (13, '2'), (14, '4'), (14, '2'), (15, '3'), (15, '2'), (16, '4'), (16, '2'), (17, '4'), (17, '2'), (18, '6'), (19, '6'), (20, '6'), (21, '6'), (22, '6'), (23, '6'), (24, '6'), (25, '6'), (26, '6'), (27, '6'), (28, '6'), (29, '6'), (30, '6'), (31, '4'), (32, '7'), (33, '4'), (33, '2'), (34, '4'), (34, '2'), (35, '4'), (35, '8'), (35, '1'), (36, '6'), (37, '4'), (37, '2'), (38, '9'), (38, '2'), (39, '7'), (40, '4'), (40, '6'), (40, '2'), (41, '4'), (41, '6'), (41, '2'), (42, '4'), (42, '2'), (43, '5'), (44, '5'), (44, '2'), (45, '5'), (45, '2'), (46, '5'), (46, '2'), (47, '5'), (47, '2'), (48, '1'), (48, '5'), (48, '2'), (49, '1'), (50, '1'), (50, '8'), (50, '2'), (51, '5'), (51, '6'), (52, '10');
