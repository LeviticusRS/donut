# donut

## Conventions

- In structures, fields should be abbreviated when appropriate (when the name would shadow a package name). Example: using `msg` instead of `message` as a field name.
- Errors that are not used by outside packages should be inlined into the functions that are used by. Exceptions to this can be made for situations where an error is reused privately by the package.
