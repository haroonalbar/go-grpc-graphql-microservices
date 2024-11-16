# This Dockerfile creates a custom PostgreSQL container that will:
# 1. Use PostgreSQL 10.3 as its base
# 2. Initialize the database with whatever SQL commands are in your `up.sql` file
# 3. Run PostgreSQL server when started
#
# This is commonly used when you want to create a PostgreSQL container with a pre-initialized database structure or data.
FROM postgres:10.3
#    - This is the base image declaration
#    - It pulls the official PostgreSQL image version 10.3 from Docker Hub
#    - This will be used as the starting point for your custom image

COPY up.sql /docker-entrypoint-initdb.d/1.sql
#    - This command copies a file named `up.sql` from your local build context
#    - It places it in the container at `/docker-entrypoint-initdb.d/1.sql`
#    - The `/docker-entrypoint-initdb.d/` directory is special in PostgreSQL Docker images
#    - Any `.sql` files in this directory will be automatically executed when the container first starts up
#    - Files are executed in alphabetical order, hence naming it `1.sql` if you need to control execution order

CMD ["postgres"]
#    - This is the default command that will run when the container starts
#    - It starts the PostgreSQL server process
#    - Note: This line is actually unnecessary here because the base postgres image already includes this CMD
