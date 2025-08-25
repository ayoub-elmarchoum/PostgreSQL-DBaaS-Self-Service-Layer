#!/bin/bash
set -e

# Check that all modules declared in .tf files start with 'sxyz'
grep -o 'module "[^"]*"' *.tf | while read -r line; do
  name=$(echo "$line" | sed -E 's/module "([^"]*)"/\1/')
  if [[ "$name" != dbaas-pg-db* ]]; then
    echo "Error: Module name '$name' does not start with 'sxyz'."
    exit 1
  fi
done

# Proceed with terraform init and validate
terraform init -input=false
terraform validate
terraform plan -out=tfplan

echo "Tests completed successfully."
