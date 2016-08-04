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
