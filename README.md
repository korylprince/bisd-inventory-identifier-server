# Info

This is the back end server for [bisd-inventory-identifier-client](https://github.com/korylprince/bisd-inventory-identifier-client), a Chrome extension to show information about BISD chromebooks.

# Install

```
go get github.com/korylprince/bisd-inventory-identifier-server
```

Create a MySQL database with `model.sql`. (This matches [pyInventory](https://github.com/korylprince/pyInventory).)

# Configuration

    INVENTORY_OUATHJSONPATH="/path/to/file.json"
    INVENTORY_OUATHIMPERSONATEUSER="user@domain.tld"
    INVENTORY_SQLDRIVER="mysql"
    INVENTORY_SQLDSN="username:password@tcp(server:3306)/database?parseTime=true"
    INVENTORY_LISTENADDR=":8080"
    INVENTORY_PREFIX="/inventory" #URL prefix
