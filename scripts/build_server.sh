#!/bin/bash
set -e
cp -r locales/ internal/mail/ && go build