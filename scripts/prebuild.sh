#!/bin/bash
set -e
./scripts/gen_version.sh
cp -r locales/ internal/mail/
