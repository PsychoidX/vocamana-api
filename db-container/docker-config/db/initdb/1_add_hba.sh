#!/bin/bash
echo "hostnossl $POSTGRES_DB $POSTGRES_USER 192.168.100.3/0 password" >> /var/lib/postgresql/data/pg_hba.conf
