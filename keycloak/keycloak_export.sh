docker exec frans_dev_keycloak mkdir -p /opt/keycloak/data/realm_import
docker exec frans_dev_keycloak /opt/keycloak/bin/kc.sh export --optimized --file /opt/keycloak/data/realm_import/dev-realm.json --realm dev
docker cp frans_dev_keycloak:/opt/keycloak/data/realm_import ./keycloak/