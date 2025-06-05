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
                                        serve_lead_name TEXT DEFAULT '',
                                        serve_lead_email TEXT DEFAULT '',
                                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                        status status NOT NULL DEFAULT 'open',
                                        website TEXT NOT NULL DEFAULT '',
                                        ages TEXT NOT NULL DEFAULT 'All Ages',
                                        leads JSONB DEFAULT '{}'::jsonb
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

INSERT INTO projects (google_id, title, description, time, project_date, max_capacity, area, latitude, longitude, serve_lead_id, serve_lead_name, serve_lead_email, location_address, website, ages) VALUES (1, 'Aging Resources Douglas County (Home 01)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (2, 'Aging Resources Douglas County (Home 02)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (3, 'Aging Resources Douglas County (Home 03)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (4, 'Aging Resources Douglas County (Home 04)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (5, 'Aging Resources Douglas County (Home 05)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (6, 'Aging Resources Douglas County (Home 06)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (7, 'Aging Resources Douglas County (Home 07)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (8, 'Aging Resources Douglas County (Home 08)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (9, 'Aging Resources Douglas County (Home 09)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (10, 'Aging Resources Douglas County (Home 10)', 'Aging Resources of Douglas County (ARDC) connects seniors with resources that help them age independently. This Serve Day, volunteers are needed to assist seniors with yard work at their homes. Tasks may include trimming branches, planting bulbs, spreading mulch, raking pine needles, etc. 

The location will be in Douglas County from 9:00 a.m. to 12:00 p.m, with the specific address provided closer to Serve Day. All ages are welcome, but children must be accompanied by an adult. Volunteers younger than middle school age will not be counted in the volunteer numbers. 

All adults over the age of 18 must complete and sign the ARDC Volunteer Application and Waiver at the link below:

http://www.agingresourcesdougco.org/service-saturdays.html 
', '9am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Douglas County', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', 'TBD. Communicated week prior to project.', 'www.agingresourcesdougco.org', 'All Ages'), (11, 'Alternatives Pregnancy Center', 'Alternatives Pregnancy Center cares for Denver-area women and men in pregnancy-related crises and offers them a meaningful alternative to abortion.

Help us create a “baby shower in a bag” for every mom who chooses life by dropping off new, unopened baby items and gift cards to the Alternatives Pregnancy Center in Englewood between 10:00 a.m. and 11:00 a.m. 
 
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
·  pacifiers', '10 - 11am', '2025-07-12 00:00:00+00:00', 30, 'Englewood', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '23 Inverness Way E. Suite 101B
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
·  gluten-free items', '9 - 10am', '2025-07-12 00:00:00+00:00', 30, 'Highlands Ranch', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '213 W. County Line Road
Highlands Ranch, CO 80129', 'www.backpacksociety.org', 'All Ages'), (13, 'Bin Blessed', 'Bin Blessed supports nonprofit There with Care in providing a wide range of meaningful services to families with children facing medical crises. This Serve Day, volunteers are needed to assemble movie cups containing candy and popcorn, and to decorate grocery bags that are used to deliver grocery items directly to families doorstep with a critically ill child. 
 
This project will take place in Auditorium #2 at Journey Church in Castle Pines from 10:00 a.m. to 12:00 p.m. Reserved for families with small children so they can serve together! 
 
Sign up to bring candy and popcorn using the Signup Genius link below:

https://www.signupgenius.com/go/5080548AFAD2AABF85-56498915-service', '10am - 12pm', '2025-07-12 00:00:00+00:00', 30, 'Journey Church (Castle Pines)', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '9009 Clydesdale Road
Castle Pines, CO 80108', 'www.binblessed.com', 'Families with Small Children ONLY'), (14, 'Bridge of Hope', 'Bridge of Hope is a nonprofit committed to ending homelessness for single mothers and their children in the greater Denver area.
For Serve Day, we are collecting backpacks and school supplies for school-aged children in need. 

THIS IS A GREAT PROJECT FOR THOSE THAT CAN''T SERVE ON SERVE DAY, or want to do another project on top of Serve Day. We will be collecting backpacks along with the supplies on Sunday, July 27th, at the Journey Castle Pines location between the times of 9 am and 11:15 am. The Serve Day Lead (Nicole Hawkins) will email the supply list to those who sign up as each kid goes to a different school and grade.', '9 - 11:15 am', '2025-07-27 00:00:00+00:00', 30, 'Journey Church (Castle Pines)', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '9009 Clydesdale Road
Castle Pines, CO 80108', 'www.greaterdenver.bridgeofhopeinc.org', 'All Ages'), (15, 'Colorado Feeding Kids (Shift 1)', 'Colorado Feeding Kids is a nonprofit whose mission is to provide nutritious food to impoverished people in Colorado and around the world. Serve Day volunteers will pack food using measuring cups, bags, and a heat press. 
 
This project will take place at 14107 E. Exposition Ave., Aurora, CO 80012, from 11:00 a.m. to 1:00 p.m. Children ages 7 and up are welcome to participate alongside a parent.', '11am - 1pm ', '2025-07-12 00:00:00+00:00', 80, 'Aurora', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '14107 E. Exposition Ave.  
Aurora, CO 80012', 'www.cofeedingkids.org', '7 Years and Older'), (16, 'Colorado Feeding Kids (Shift 2)', 'Colorado Feeding Kids is a nonprofit whose mission is to provide nutritious food to impoverished people in Colorado and around the world. Serve Day volunteers will pack food using measuring cups, bags, and a heat press. 
 
This project will take place at 14107 E. Exposition Ave., Aurora, CO 80012, from 1:30 p.m. to 3:30 p.m. Children ages 7 and up are welcome to participate alongside a parent.', '1:30 - 3:30pm', '2025-07-12 00:00:00+00:00', 80, 'Aurora', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '14107 E. Exposition Ave.  
Aurora, CO 80012', 'www.cofeedingkids.org', '7 Years and Older'), (17, 'Denver Rescue Mission - 48th Ave Center (Shift 1)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 7:15 a.m. to 9:00 a.m. 

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '7:15 - 9:00am', '2025-07-12 00:00:00+00:00', 6, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '14 Years and Older'), (18, 'Denver Rescue Mission - 48th Ave Center (Shift 2)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will serve the homeless and those in transition through participating in the lunch meal service (including food preparation and distribution) from 11:15 a.m. to 1:00 p.m. 

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '11:15am - 1pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '14 Years and Older'), (19, 'Denver Rescue Mission - 48th Ave Center (Shift 4)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will serve the homeless and those in transition through participating in the dinner meal service (including food preparation and distribution) from 4:45 p.m. to 6:30 p.m. 

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '4:45 - 6:30pm', '2025-07-12 00:00:00+00:00', 8, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '14 Years and Older'), (20, 'Denver Rescue Mission - Lawrence St. Center (Shift 1)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 6:45 a.m. to 8:45 a.m. 

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '6:45 - 8:45am', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '14 Years and Older'), (21, 'Denver Rescue Mission - Lawrence St. Center (Shift 2)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will help with dishwashing following the breakfast service from 7:30 a.m. to 9:00 a.m. 
 
Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 18 or older.', '7:30 - 9:00am', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '18 Years and Older'), (22, 'Denver Rescue Mission - Lawrence St. Center (Shift 3)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will help with dishwashing following the lunch service from 12:30 p.m. to 2:00 p.m. 

Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 18 or older.', '12:30 - 2pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '18 Years and Older'), (23, 'Denver Rescue Mission - Lawrence St. Center (Shift 4)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will help prepare the dinner service from 2:30 p.m. to 3:30 p.m. 
 
Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '2:30 - 3:30pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '14 Years and Older'), (24, 'Denver Rescue Mission - Lawrence St. Center (Shift 5)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will serve the homeless and those in transition through participating in the dinner meal service (including food preparation and distribution) from 4:45 p.m. to 6:30 p.m. 
 
Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 14 or older; those aged 14 to 17 must be accompanied by an adult.', '4:45 - 6:30pm', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '14 Years and Older'), (25, 'Denver Rescue Mission - Lawrence St. Center (Shift 6)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will help with dishwashing following the dinner service from 5:30 p.m. to 7:00 p.m. 
 
Denver Rescue Mission is located at 2222 Lawrence St., Denver, CO 80205. Volunteers must be 18 or older.', '5:30 - 7:00pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '2222 Lawrence St.
Denver, CO 80205', 'www.denverrescuemission.org', '18 Years and Older'), (26, 'Denver Rescue Mission - 48th Ave Center (Shift 3)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will help with dishwashing following the lunch service from 12:00 p.m. to 1:30 p.m. 

Denver Rescue Mission is located at 4600 E 48th Ave., Denver, CO 80216. Volunteers must be 18 or older.', '12 - 1:30pm', '2025-07-12 00:00:00+00:00', 2, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '4600 E 48th Ave.
Denver, CO 80216', 'www.denverrescuemission.org', '18 Years and Older'), (27, 'Denver Rescue Mission - The Crossing (Shift 1)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will serve the homeless and those in transition through participating in the breakfast meal service (including food preparation and distribution) from 6:00 a.m. to 7:15 a.m. 

Denver Rescue Mission is located at 6090 Smith Rd., Denver, CO 80216. Volunteers must be 12 or older; those aged 12 to 17 must be accompanied by an adult.', '6 - 7:15am', '2025-07-12 00:00:00+00:00', 3, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '6090 Smith Road
Denver, CO 80216', 'www.denverrescuemission.org', '12 Years and Older'), (28, 'Denver Rescue Mission - The Crossing (Shift 2)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will assist with afternoon kitchen duties from 2:00 p.m. to 4:00 p.m. 

Denver Rescue Mission is located at 6090 Smith Rd., Denver, CO 80216. Volunteers must be 12 or older; those aged 12 to 17 must be accompanied by an adult.', '2 - 4pm', '2025-07-12 00:00:00+00:00', 3, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '6090 Smith Road
Denver, CO 80216', 'www.denverrescuemission.org', '12 Years and Older'), (29, 'Denver Rescue Mission - The Crossing (Shift 3)', 'Denver Rescue Mission is committed to helping people who are experiencing homelessness and addiction change their lives. 

Volunteers will serve the homeless and those in transition through participating in the dinner meal service (including food preparation and distribution) from 5:00 p.m. to 7:00 p.m. 

Denver Rescue Mission is located at 6090 Smith Rd., Denver, CO 80216. Volunteers must be 12 or older; those aged 12 to 17 must be accompanied by an adult.', '5 - 7pm', '2025-07-12 00:00:00+00:00', 6, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Allison Adams', 'smada3@comcast.net', '6090 Smith Road
Denver, CO 80216', 'www.denverrescuemission.org', '12 Years and Older'), (30, 'Denver Street School (Project 1)', 'The Denver Street School is a high school program that brings the love of Christ to students who are struggling academically. It provides a quality education as well as support to help kids thrive. 
 
Volunteers will assist with cleaning and organizing the teacher resource room from 9:00 a.m. to 12:00 p.m.
 
The Denver Street School is located at 1380 Ammons St., Lakewood, CO 80214. Children are welcome to participate alongside a parent.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 5, 'Lakewood', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', '1380 Ammons St 
Lakewood, CO 80214', 'www.denverstreetschool.org', 'All Ages'), (31, 'Denver Street School (Project 2)', 'The Denver Street School is a high school program that brings the love of Christ to students who are struggling academically. It provides a quality education as well as support to help kids thrive. 
 
Volunteers will assist with painting the teacher lounge from 9:00 a.m. to 12:00 p.m.
 
The Denver Street School is located at 1380 Ammons St., Lakewood, CO 80214. This project is for those 18 years and older.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 5, 'Lakewood', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', '1380 Ammons St 
Lakewood, CO 80214', 'www.denverstreetschool.org', '18 Years and Older'), (32, 'Food Bank of the Rockies (Shift 1)', 'Food Bank of the Rockies provides food and other necessities to people facing hunger in Colorado and Wyoming. 
 
Volunteers will report to the main warehouse located at 10700 E 45th Ave., Denver, CO 80239, where they will inspect, clean, sort, and repack perishable and nonperishable food items for distribution to Hunger Relief Partners from 8:45 a.m. to 12:00 p.m. All group participants must be signed up two weeks prior to Serve Day, by June 28, 2025. Children ages 10 and up are welcome to participate alongside an adult chaperone (see chaperone requirements below).

Food Bank of the Rockies requires each volunteer to register and sign a waiver ahead of volunteering and will allow them to track their volunteer hours with Food Bank of the Rockies. 

To register, visit https://vhub.at/1S4HE8D. Click on the date and the green “Sign up” button next to the event on the page. Then you will be prompted to either log in or create an account. You do not need a join code and can skip this step. Once you have created an account and registered for the event, you will receive a confirmation email. 

If your group includes 10- and 11-year-olds, please have 1 chaperone for every 3 youths. 
If your group has 12- to 14-year-olds, please have 1 chaperone for every 4 youths. 
If your group has 15- to 17-year-olds, please have 1 chaperone for every 5 youths. 
Peers 18 years of age do not count as chaperones. 

If you have any questions, please reach out to volunteer@foodbankrockies.org.', '8:45am - 12pm', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '10700 E. 45th Ave.  
Denver, CO 80239', 'www.foodbankrockies.org', '10 Years and Older'), (33, 'Food Bank of the Rockies (Shift 2)', 'Food Bank of the Rockies provides food and other necessities to people facing hunger in Colorado and Wyoming. 
 
Volunteers will report to the main warehouse located at 10700 E 45th Ave., Denver, CO 80239, where they will inspect, clean, sort, and repack perishable and nonperishable food items for distribution to Hunger Relief Partners from 1:00 p.m. to 4:00 p.m. All group participants must be signed up two weeks prior to Serve Day, by June 28, 2025. Children ages 10 and up are welcome to participate alongside an adult chaperone (see chaperone requirements below).

Food Bank of the Rockies requires each volunteer to register and sign a waiver ahead of volunteering and will allow them to track their volunteer hours with Food Bank of the Rockies. 

To register, visit https://vhub.at/1S4HE8D. Click on the date and the green “Sign up” button next to the event on the page. Then you will be prompted to either log in or create an account. You do not need a join code and can skip this step. Once you have created an account and registered for the event, you will receive a confirmation email. 

If your group includes 10- and 11-year-olds, please have 1 chaperone for every 3 youths. 
If your group has 12- to 14-year-olds, please have 1 chaperone for every 4 youths. 
If your group has 15- to 17-year-olds, please have 1 chaperone for every 5 youths. 
Peers 18 years of age do not count as chaperones. 

If you have any questions, please reach out to volunteer@foodbankrockies.org.', '1 - 4pm', '2025-07-12 00:00:00+00:00', 24, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '10700 E. 45th Ave.  
Denver, CO 80239', 'www.foodbankrockies.org', '10 Years and Older'), (34, 'Help and Hope Center', 'Help & Hope Center supports Douglas and Elbert County residents in financial distress, allowing them to navigate troublesome times with dignity. 

Volunteers will pack Strive to Thrive bags, sort donations in receiving area, perform light cleaning projects, and assist with outdoor work such as weeding and general beautification. Volunteering will take place from 9:00 a.m. to 11:00 a.m. at 1638 Park St., Castle Rock, CO 80109.
 
Be advised that work may be performed outside and/or dirty, so participants should wear weather-appropriate wear clothing and closed-toe shoes. Volunteers must be at least 13 years old and accompanied by an adult.', '9 - 11am', '2025-07-12 00:00:00+00:00', 20, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '1638 Park St.
Castle Rock, CO 80109', 'www.helpandhopecenter.org', '13 Years and Older'), (35, 'Hope''s Promise (Project 1)', 'Hope’s Promise strives to strengthen families through ethical, Christ-centered foster care, adoption, and global orphan care. 
 
Volunteers will make welcoming meals for foster children and their new families. Drop off welcoming meals at Hope’s Promise between 10:00 a.m. and 12:00 p.m. at 1585 S Perry St., #E, Castle Rock, CO 80104.

Make a freezer meal for a family of 4-5 people. Label meal with reheating instructions and a list ingredients. Bring frozen cookie dough for dessert, too. Any volunteers who would like a tour of Hope’s Promise can arrive early 9:45 a.m.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 20, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '1585 S Perry St # E
Castle Rock, CO 80104', 'www.hopespromise.com', 'All Ages'), (36, 'Hope''s Promise (Project 2)', 'Hope’s Promise strives to strengthen families through ethical, Christ-centered foster care, adoption, and global orphan care. 
 
Volunteers will help prepare for an upcoming golf tournament by helping to assemble silent auction items and baskets along with creating signage and banquet decor.
 
Volunteering will take place at 1585 S Perry St., #E, Castle Rock, CO 80104 from 10:00 a.m. to 12:00 p.m.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 7, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '1585 S Perry St # E
Castle Rock, CO 80104', 'www.hopespromise.com', '10 Years and Older'), (37, 'Hope''s Promise (Project 3)', 'Hope’s Promise strives to strengthen families through ethical, Christ-centered foster care, adoption, and global orphan care. 
 
Volunteers will write encouraging, uplifting notes to parents, caregivers, and foster parents both internationally and in Colorado. Please bring a package of stationary or pretty note cards for this project. 
 
Volunteering will take place at 1585 S Perry St., #E, Castle Rock, CO 80104 from 10:00 a.m. to 12:00 p.m.', '10am - 12pm', '2025-07-12 00:00:00+00:00', 15, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '1585 S Perry St # E
Castle Rock, CO 80104', 'www.hopespromise.com', 'All Ages'), (38, 'His Love Fellowship Church', 'His Love Fellowship is a Bible-based church in west Denver.
 
Volunteers will help stain railroad ties around the trees and paint over graffiti. His Love Fellowship is located at 910 Kalamath St., Denver, CO 80204.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 10, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', '910 Kalamath St.
Denver, CO 80204', 'www.hislovefellowship.com', '12 Years and Older'), (39, 'SECOR (Shift 1)', 'SECOR cares for people faced with suburban poverty by offering a Free Food Market where guests can shop for food and other necessities. Most guests are able to go home with enough food for 10-14 days.

This project will involve both market and warehouse service at SECOR, located at 17151 Pine Ln., Parker, CO 80134, from 8:45 a.m. to 11:00 a.m. 

Market volunteers will help guests get started after they check in with the front desk, man various stations in the market as guests shop, and manage the cart return. Warehouse volunteers will help sort items that come in on the truck, restock market shelves, and help clean the market and warehouse toward the end of the shift. They may also assist with packing boxes, sorting donations, and assembling items. The first shift of volunteers will show the second shift the ropes, as SECOR will not have enough staff to handle that midday while guests are shopping.

Volunteers must be 8 years old and accompanied by an adult (groups of minors need to have a 1 adult to 6 minors ratio). Closed-toe shoes are required. All volunteers must fill out a liability waiver: https://forms.gle/ymctoRUt2iPfjSYQ8.

Please meet in the lobby of the building on the east end of the parking lot. Bring a lock if you want to lock up any valuables. Lastly, please stay home if you are feeling ill.', '8:45 - 11am', '2025-07-12 00:00:00+00:00', 10, 'Parker', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '17151 Pine Lane
Parker, CO 80134', 'www.secorcares.com', '8 Years and Older'), (40, 'SECOR (Shift 2)', 'SECOR cares for people faced with suburban poverty by offering a Free Food Market where guests can shop for food and other necessities. Most guests are able to go home with enough food for 10-14 days.

This project will involve both market and warehouse service at SECOR, located at 17151 Pine Ln., Parker, CO 80134, from 10:45 a.m. to 1:15 p.m. 

Market volunteers will help guests get started after they check in with the front desk, man various stations in the market as guests shop, and manage the cart return. Warehouse volunteers will help sort items that come in on the truck, restock market shelves, and help clean the market and warehouse toward the end of the shift. They may also assist with packing boxes, sorting donations, and assembling items. The first shift of volunteers will show the second shift the ropes, as SECOR will not have enough staff to handle that midday while guests are shopping.

Volunteers must be 8 years old and accompanied by an adult (groups of minors need to have a 1 adult to 6 minors ratio). Closed-toe shoes are required. All volunteers must fill out a liability waiver: https://forms.gle/ymctoRUt2iPfjSYQ8.

Please meet in the lobby of the building on the east end of the parking lot. Bring a lock if you want to lock up any valuables. Lastly, please stay home if you are feeling ill.', '10:45am - 1:15pm', '2025-07-12 00:00:00+00:00', 10, 'Parker', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '17151 Pine Lane
Parker, CO 80134', 'www.secorcares.com', '8 Years and Older'), (41, 'Sweet Dream In A Bag', 'Sweet Dream in a Bag has provided cozy bedding to children in need since 2010. For this project, volunteers will assemble bags with bedding to be distributed to children in need, both locally and internationally.
 
Event will run from 9:00 a.m. to 11:00 a.m. at 5933 S. Fairfield St., Littleton, CO 80120. Kids of all ages are welcome; those under age 5 will be more like helpers to their parents.', '9 - 11am', '2025-07-12 00:00:00+00:00', 20, 'Littleton', 39.491482, -104.874878, 'example-user-123', 'Jenn Mell', 'jennmell@me.com', '5933 S. Fairfield Street
Littleton, CO 80120', 'www.sweetdreaminabag.org', 'All Ages'), (42, 'Vitalant Blood Bank', 'This Serve Day, consider participating in our blood drive! 
 
The blood drive will take place in the upstairs atrium at Journey Church in Castle Pines between 9:00 a.m. and 1:00 p.m. Blood donors must be 16 years or older and weigh at least 110 pounds. 
 
To sign up, click here: https://donors.vitalant.org/dwp/portal/dwa/appointment/guest/phl/timeSlotsExtr?token=YejHtw9ohlDWKytYWH8lXWexFLCpLg9%2Bu7xOmseb9Qs%3D

To confirm your eligibility to donate blood, see https://www.vitalant.org/eligibility', '9am - 1pm', '2025-07-12 00:00:00+00:00', 56, 'Journey Church (Castle Pines)', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '9009 Clydesdale Road
Castle Pines, CO 80108', 'www.vitalant.org', '16 Years and Older'), (43, 'WeSoThey (Shift 1)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally. This Serve Day, volunteers are needed to help with WeSoThey’s annual fundraising summer garage sale.

For this project, volunteers will help with organizing, setting up, assisting morning shoppers from 7:30 a.m. to 10:30 a.m. at 800 Third St., Castle Rock, CO 80104. This is a family-friendly activity suitable for kids ages 5 and up and their parents.', '7:30 - 10:30am', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '880 Third Street 
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (44, 'WeSoThey (Shift 2)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally. This Serve Day, volunteers are needed to help with WeSoThey’s annual fundraising summer garage sale.

For this project, volunteers will help with organizing tables and assisting shoppers from 10:30 a.m. to 1:00 p.m. at 800 Third St., Castle Rock, CO 80104. This is a family-friendly activity suitable for kids ages 5 and up and their parents.', '10:30am - 1pm', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '880 Third Street 
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (45, 'WeSoThey (Shift 3)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally. This Serve Day, volunteers are needed to help with WeSoThey’s annual fundraising summer garage sale.

For this project, volunteers will help with organizing tables and assisting shoppers from 1:00 p.m. to 4:00 p.m. at 800 Third St., Castle Rock, CO 80104. This is a family-friendly activity suitable for kids ages 5 and up and their parents.', '1 - 4pm', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '880 Third Street 
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (46, 'WeSoThey (Shift 4)', 'WeSoThey supports orphanages in Uganda and Mexico and families adopting internationally. This Serve Day, volunteers are needed to help with WeSoThey’s annual fundraising summer garage sale.

For this project, volunteers will help with tear-down from 4:00 p.m. to 6:00 p.m. at 800 Third St., Castle Rock, CO 80104. Tear-down duties are reserved for adults who are able to lift heavy items.', '4 - 6pm', '2025-07-12 00:00:00+00:00', 10, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '880 Third Street 
Castle Rock, Co', 'www.wesothey.com', '5 Years and Older'), (47, 'Leman Academy (Journey Parker)', 'Volunteers will have the opportunity to love on the school that hosts Journey’s Parker location by cleaning up weeds and trash around the school, tidying the curb areas (sweeping rocks and picking up any trash), sweeping wood chips from the southside sidewalk back into the play areas, and disposing of any broken/unusable cones. 
 
Meet at Lehman Academy, 19560 Stroh Rd., Parker, CO 80134, from 9:00 a.m. to 12:00 p.m. All ages are welcome; kids must be accompanied by an adult.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 20, 'Parker', 39.491482, -104.874878, 'example-user-123', 'Cindy Kavanaugh', 'cindykavanaugh@journeycolorado.com', '19560 Stroh Rd.
Parker, CO 80134', 'N/A', 'All Ages'), (48, 'Barksdale Family Home Project', 'The Barksdale Family home could use some volunteers to assist with much-needed yard work, including cutting brush and scrub oak and putting them in piles. 
 
This activity requires the ability to perform physical labor and is for adults only. The project will take place from 8:00 a.m. to 11:00 p.m. Please wear long sleeves, a hat, gloves, and closed-toe shoes. Please also bring a chainsaw or any tools that can cut brush.

The Serve Day Lead will provide volunteers with the location address as the date gets closer.', '8 - 11am', '2025-07-12 00:00:00+00:00', 20, 'Castle Rock', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', 'Address will be provided separately for those registered', 'N/A', '8 Years and Older'), (49, 'Open Door Ministries', 'Open Door Ministry’s Cornerstone Home is a four-bedroom shelter for disabled men in Denver.

This project involves doing a deep clean of Cornerstone’s main living areas, including the living room, kitchen, and hallways. Volunteers will also do some light landscape work such as planting pots in front of the home and weeding the mulch beds.

Volunteering will take place at 2754 Stout St., Denver, CO 80205, from 9:00 a.m. to 12:00 p.m. Please wear long sleeves and bring gloves, a hat, sunscreen, water, and cleaning supplies. All ages are welcome.', '9am - 12pm', '2025-07-12 00:00:00+00:00', 15, 'Denver', 39.491482, -104.874878, 'example-user-123', 'Nicole Hawkins', 'nhawkins10322@gmail.com', '2754 Stout St.  
Denver, CO 80205', 'www.odmdenver.org', 'All Ages'), (50, 'Parker Task Force', 'The Parker Task Force assists individuals and families in the community with everything they need to become self-sufficient. 
 
Serve Day volunteers are needed to assist with the food bank facility by helping shoppers from 8:30 a.m. to 1:00 p.m. The food bank is located at 19105 Longs Way, Parker, CO 80134, and the shift will begin with a tour of the facility along with some training. This project is for those 18 years and older.', '8:30am - 1pm', '2025-07-12 00:00:00+00:00', 4, 'Parker', 39.491482, -104.874878, 'example-user-123', 'Cindy Kavanaugh', 'cindykavanaugh@journeycolorado.com', '19105 Longs Way 
Parker, CO 80134', 'www.parkertaskforce.org', '18 Years and Older'), (51, 'Affinity Ranch (Project 1)', 'Affinity Ranch serves the community, including veterans and those with differing abilities, by offering healing with animals. They currently care for 18 horses, making sure they are healthy and happy. Volunteers will work with ranch staff to build and donate a new hay feeder for the horses.

This activity will take place at 11892 Hilltop Rd., Parker, CO 80134, from 9:00 a.m. to 1:00 p.m. Affinity Ranch will provide the building supplies, but volunteers should bring gloves and safety glasses, and they should wear closed-toe work shoes or boots as well as long sleeves and a hat. 
 
A waiver must be signed ahead of time:

https://prayinghandsranch.org/ranchapps/PublicForms/VolunteerGroupRegister.php?GrpName=JOURNEYCHURCH', '9am - 1pm', '2025-07-12 00:00:00+00:00', 5, 'Parker', 39.491482, -104.874878, 'example-user-123', 'Cindy Kavanaugh', 'cindykavanaugh@journeycolorado.com', '11892 Hilltop Road
Parker, 80134', 'www.affinityranch.org', '18 Years and Older'), (52, 'Affinity Ranch (Project 2)', 'Affinity Ranch serves the community, including veterans and those with differing abilities, by offering healing with animals. At the ranch, the work never ends. Volunteers will help with general chores and projects around the ranch, including possibly feeding animals and cleaning stalls. 

This activity will take place at 11892 Hilltop Rd., Parker, CO 80134, from 9:00 a.m. to 1:00 p.m. Volunteers should wear closed-toe shoes/boots and work clothes that can get dirty. Long sleeves and a hat are suggested. 
 
A waiver must be signed ahead of time:

https://prayinghandsranch.org/ranchapps/PublicForms/VolunteerGroupRegister.php?GrpName=JOURNEYCHURCH', '9am - 1pm', '2025-07-12 00:00:00+00:00', 10, 'Parker', 39.491482, -104.874878, 'example-user-123', 'Cindy Kavanaugh', 'cindykavanaugh@journeycolorado.com', '11892 Hilltop Road
Parker, 80134', 'www.affinityranch.org', '18 Years and Older'), (53, 'DIY Project - Love Your Neighbor', 'Are all the Serve Day projects taken, or are the ones that work best for your situation on July 12th unavailable? The great thing is that you can still serve. Just "DIY" the project! Come up with creative ways to serve someone around you – AKA, love your neighbor. You come up with the idea and the timing, but here are some ideas. Deliver cookies/treats/lunch to your local police station or fire station. Mow an elderly neighbor’s yard, pick weeds, wash their car, or run errands for them. Have a lemonade or other type of stand and donate the proceeds to one of the Serve Day organizations. Go on a walk or drive and pray together as a family or small group for households and/or individuals. Deliver dinner or treats to a neighbor or someone in need. Pick up trash around a neighborhood, school, or park. Serve root beer floats to all your neighbors on your street. 

You can do it anytime, anywhere, and in any way. Just "DIY" it!', 'Anytime', '2025-07-12 00:00:00+00:00', 100, 'Anywhere', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', 'Anywhere', 'N/A', 'All Ages'), (54, 'Vacation Bible School Prep (Shift 1)', 'Vacation Bible School (VBS) is one of Journey''s largest Kids Ministry''s events of the year. There are over 300 kids signed up for this week-long experience at the Journey Castle Pines location.

This project will involve a lot of different ways to serve in preparation for VBS. There will be opportunities to decorate, clean, organize, sort, and simply be available in any way to help transform Journey into an environment for kids to take their next steps with Jesus.

Shift 1 will be from 10 am to 12 pm at Journey Castle Pines. Volunteers will need to be 13 years and older.
', '10am - 12pm', '2025-07-12 00:00:00+00:00', 25, 'Journey Church (Castle Pines)', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '9009 Clydesdale Road
Castle Pines, CO 80108', 'N/A', '12 Years and Older'), (55, 'Vacation Bible School Prep (Shift 2)', 'Vacation Bible School (VBS) is one of Journey's largest Kids Ministry's events of the year. There are over 300 kids signed up for this week-long experience at the Journey Castle Pines location.

This project will involve a lot of different ways to serve in preparation for VBS. There will be opportunities to decorate, clean, organize, sort, and simply be available in any way to help transform Journey into an environment for kids to take their next steps with Jesus.

Shift 2 will be from 12 pm to 2 pm at Journey Castle Pines. Volunteers will need to be 13 years and older.
', '12pm - 2pm', '2025-07-12 00:00:00+00:00', 25, 'Journey Church (Castle Pines)', 39.491482, -104.874878, 'example-user-123', 'Nicole Cecil', 'nicoleacecil@gmail.com', '9009 Clydesdale Road
Castle Pines, CO 80108', 'N/A', '12 Years and Older');
INSERT INTO types (id, type)  VALUES (7, 'Cleaning/Organizing'), (8, 'Painting'), (10, 'Arts & Crafts'), (1, 'Yard Work'), (2, 'Kid Friendly'), (9, 'Minor Indoor Repairs'), (11, 'Construction'), (3, 'Collection/Drop Off'), (4, 'Sorting/Assembly'), (5, 'Community Outreach'), (6, 'Food Prep/Distribution');
INSERT INTO project_types (project_id, type_id)  VALUES (1, '1'), (1, '2'), (2, '1'), (2, '2'), (3, '1'), (3, '2'), (4, '1'), (4, '2'), (5, '1'), (5, '2'), (6, '1'), (6, '2'), (7, '1'), (7, '2'), (8, '1'), (8, '2'), (9, '1'), (9, '2'), (10, '1'), (10, '2'), (11, '3'), (11, '2'), (12, '3'), (12, '2'), (13, '4'), (13, '5'), (13, '2'), (14, '3'), (14, '2'), (15, '4'), (15, '2'), (16, '4'), (16, '2'), (17, '6'), (18, '6'), (19, '6'), (20, '6'), (21, '7'), (22, '7'), (23, '6'), (24, '6'), (25, '7'), (26, '7'), (27, '6'), (28, '6'), (29, '6'), (30, '7'), (30, '2'), (31, '8'), (32, '4'), (32, '2'), (33, '4'), (33, '2'), (34, '4'), (34, '9'), (34, '1'), (35, '6'), (36, '4'), (36, '2'), (37, '10'), (37, '2'), (38, '8'), (39, '4'), (39, '6'), (39, '2'), (40, '4'), (40, '6'), (40, '2'), (41, '4'), (41, '2'), (42, '5'), (43, '5'), (43, '7'), (43, '2'), (44, '5'), (44, '7'), (44, '2'), (45, '5'), (45, '7'), (45, '2'), (46, '5'), (46, '7'), (46, '2'), (47, '1'), (47, '5'), (47, '2'), (48, '1'), (49, '1'), (49, '7'), (49, '9'), (49, '2'), (50, '5'), (50, '6'), (51, '11'), (52, '7'), (52, '1'), (53, '2'), (54, '10'), (54, '4'), (54, '7'), (55, '10'), (55, '4'), (55, '7');
