# RIPE AS-SET updater
Quick add or delete ASN/AS-SET to your AS-SET through [RIPE RESTful API](https://www.ripe.net/manage-ips-and-asns/db/support/documentation/ripe-database-documentation/updating-objects-in-the-ripe-database/6-1-restful-api).

# Requirenments
You need to add "auth:" attribute in mntner object which protects your as-set object. "auth:" attribute must be in [MD5 format](https://www.ripe.net/manage-ips-and-asns/db/support/security/protecting-data#MD5)

# How to use

### You may just take binary and launch:
```sh
$ ./ripeRestAssetUpdate
USAGE:  ./ripeRestAssetUpdate add/delete AS/AS-SET password as-set-name
```
- add/delete - what you want to do: add or delete some as/as-set from your as-set
- AS/AS-SET - as or as-set which you want to add or delete from yours as-set
- password - is password which you added earlier to your mntner, which protects your as-set.
- as-set-name - your as-set which you want to update

example:
```sh
$ ./ripeRestAssetUpdate add AS123456 PaSSw0Rd AS-MYASSET
```
### Or you may compile the source file first of all
``` sh
$ go build ripeRestAssetUpdate.go
```
And then use obtained binary
