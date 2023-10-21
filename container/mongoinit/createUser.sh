#!/bin/bash

UserName=$MONGO_INITDB_USERNAME
Password=$MONGO_INITDB_PASSWORD
RootUserName=$MONGO_INITDB_ROOT_USERNAME
RootPassword=$MONGO_INITDB_ROOT_PASSWORD
Database=$MONGO_INITDB_DATABASE

echo "********RootUserName = " ${RootUserName}
echo "********RootPassword = " ${RootPassword}
echo "********UserName = " ${UserName}
echo "********Password = " ${Password}
echo "********Database = " ${Database}

echo "********Waiting for startup..********"

sleep 5

echo "********Started..********"


mongo -u ${MONGO_INITDB_ROOT_USERNAME} -p ${MONGO_INITDB_ROOT_PASSWORD} --authenticationDatabase admin "$MONGO_INITDB_DATABASE" <<EOF
    db.createUser({
        user: '$MONGO_INITDB_USERNAME',
        pwd: '$MONGO_INITDB_PASSWORD',
        roles: [ { role: 'readWrite', db: '$MONGO_INITDB_DATABASE' } ],
    })
EOF
