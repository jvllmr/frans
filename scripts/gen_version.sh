VERSION=$(git describe --always)

echo "package config" > internal/config/version.go
echo "" >> internal/config/version.go
echo "const FransVersion = \"$VERSION\"" >> internal/config/version.go
echo "New version: $VERSION"