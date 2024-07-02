#!/bin/sh

set -e

echo "start the app"
exec "$@"  # take all params passed to script and run it