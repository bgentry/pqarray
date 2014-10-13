# pqarray

A simple Go package for parsing PostgreSQL arrays using the `sql.Scanner`
interface.

The array splitting code was written for [PQL](https://bitbucket.org/pkg/pql)
by Chris Familoe. I extracted that portion of the code into this package and
exported it so that it can be used for any array type that you're decoding.
I've also added helper types to easily decode string and int arrays.

These types should be usable with [pq](https://github.com/lib/pq).

## TODO

* Encoding using the `database/sql` `driver.Valuer` so these types can be
  serialized back into PostgreSQL
* Multi-dimensional arrays
* Helper types for other common types, such as `time.Time`, `int64`, etc.
