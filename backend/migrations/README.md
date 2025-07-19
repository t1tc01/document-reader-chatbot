# Database Migrations

This directory contains Atlas database migrations for the document-reader-chatbot backend.

## Atlas Migration Commands

### Generate a new migration
```bash
# Generate migration from GORM models
atlas migrate diff <migration_name> \
  --dir file://migrations \
  --to "file://schema/models.sql" \
  --dev-url "postgres://postgres:password@localhost:5432/document_reader_chatbot_dev?sslmode=disable"
```

### Apply migrations
```bash
# Apply migrations to local database
atlas migrate apply --url "postgres://postgres:password@localhost:5432/document_reader_chatbot?sslmode=disable"

# Apply migrations to production
atlas migrate apply --env production
```

### Validate migrations
```bash
# Validate migration files
atlas migrate validate --dir file://migrations

# Lint migrations
atlas migrate lint --dir file://migrations --dev-url "postgres://postgres:password@localhost:5432/document_reader_chatbot_dev?sslmode=disable"
```

### Check migration status
```bash
# Check current migration status
atlas migrate status --url "postgres://postgres:password@localhost:5432/document_reader_chatbot?sslmode=disable"
```

## Migration File Naming Convention

Migration files follow the pattern: `<timestamp>_<description>.sql`

Example: `20241201120000_create_users_table.sql`

## Important Notes

1. Always test migrations on a development database first
2. Backup production database before applying migrations
3. Review generated migrations before applying them
4. Use descriptive names for migrations
5. Keep migrations small and focused 