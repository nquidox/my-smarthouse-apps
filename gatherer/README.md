# Smarthouse micropython - Gatherer

A script that gathers data from sensors via mqtt and stores it in postgresql database.

Fill in connection params for database connection.
Edit gatherer.service file if you want to use it with systemd.

In case if psycopg2 doesn't install from req.txt, install it manually with pip install psycopg2-binary
