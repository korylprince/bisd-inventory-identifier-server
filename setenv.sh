#!/bin/sh

export INVENTORY_SQLDSN="$INVENTORY_USER:$INVENTORY_PASS@tcp($INVENTORY_HOST:$INVENTORY_PORT)/$INVENTORY_DATABASE?parseTime=true"

exec $1
