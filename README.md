Description
-----------

This application parses data from the Open Street Map Protocol Buffer Format (http://wiki.openstreetmap.org/wiki/PBF_Format) and imports it into a RethinkDB database.

It has been written in GoLang and is a port of my Haskell project https://github.com/ederoyd46/OSMImport, but for RethinkDB instead of MongoDB.

Build Instructions
------------------

To build a binary, run

```
make update_deps build
```

Or to install, run

```
make update_deps install
```

Usage
-----

Three options are needed to start the import process

1. dbconnection - the host and port the database is running.
2. dbname - the name of the database you want to import into.
3. filename - the name of the file to import.


[EXAMPLE]

osmimport-go '127.0.0.1:28015' 'geo' './download/england-latest.osm.pbf'

<!-- Docker Usage
------------
Pull down the repository

```
docker pull ederoyd46/osmimport
```

Run an import, assumes you have a container called db, and have downloaded the england data from OSM in protocol buffer format

```
docker run -d --name mongo -p 27017:27017 -v $(pwd)/data:/data/db mongo
docker run -it --rm=true --link mongo:mongo -v $(pwd)/download:/data ederoyd46/osmimport 'mongo:27017' 'geo_data' '/data/england-latest.osm.pbf'
``` -->
