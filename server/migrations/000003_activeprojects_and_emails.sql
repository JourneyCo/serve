-- add column to projects to show if it is active or not
ALTER table projects ADD COLUMN active bool DEFAULT true;

-- add column to registrations to show if the two week email and one week emails have been sent
ALTER table registrations ADD COLUMN two_week_email_sent bool DEFAULT false;
ALTER table registrations ADD COLUMN one_week_email_sent bool DEFAULT false;