# Standalone example

In this `docker-compose.yml` file you find a very simple standalone example for an operational setup.

First, the instances for mysql and productservice are pulled up.

Then the migrate and seed services run to set up the DB-schema and insert some test data.
They will restart on failure, so we are sure to have the migration and seeding done, even if their first tries are too early for the mysql service.

At last, the app is started and will restart always, so just wait until migrate finishes  with exit code 0.
