# sricheck
A GoLang subresource integrity check module


# Exported Functions

## GenerateSriMap
Used to make integrity attributes from a given resource.

### Inputs
* jsSrc string - a string URL to a file, generally JS or CSS

### Returns 
* a map of hashes for that data (map[string]string). The index, or key, of the map will be a string of the hash algorithm. The value of the map will be the base64 hash for the data
* error


## CheckPageIntegrity
Used to check all of the intrgrity attributes on a page.

### Inputs
* pageUrl - a string URL to a page

### Returns
* Return a boolean value based on the integrity checks - true when all integrity attributes are valid (or not present)
and false when an integrity attribute does not check correctly.
* error

## PrintPageIntegrityCheckTable
Used to print out (to cli) a nice table for enumerating JS resources, integrity attributes, and their validity.

### Inputs
* pageUrl - a string URL to a page

### Output
Print a table summarizing all of the JS includes, their integrity attributes, and the validity of integrity attributes

### Returns
None

## CheckIntegrity
Used to check a provided integrity attribute for a provided URL.

### Inputs
* jsSrc - a string URL to a file (generally a JS or CSS file) 
* integrity - the string of the integrity attribute

### Returns
* a boolean - Returns true for valid and false for invalid
* error

# Other Todos:
[ ] Improve testing

