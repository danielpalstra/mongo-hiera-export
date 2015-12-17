# Description
This application connects to a MongoDB hiera backend, retrieves all documents and converts all hiera entries to YAML files. Purpose of the application
is migrating an existing Mongo hiera backend to a YAML backend which can be used by puppet.

# Installation
This project depends on the following libraries.
```
go get gopkg.in/yaml.v2
go get labix.org/v2/mgo
go get labix.org/v2/mgo/bson
```
Create the hieradata directory where the output will be written to.
```
mkdir hieradata
```

# Running
Set the uri to the mongo instance you want to use
```
export MONGOHQ_URL="mongodb://localhost/puppet"
```

Run the tool
```
go run hiera.go
```

# Creating test data
Create a test set to use local
## Dump from external MongoDB
```
mongodump --host $MONGOHQ_URL --port 27017 --db puppet --collection hiera --out dumps/
```

## Restore to localhost
```
mongorestore --host localhost --port 27017 --db puppet dumps/puppet/hiera.bson
```
