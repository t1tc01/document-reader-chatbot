# Atlas configuration for document-reader-chatbot backend

# Define the database environment
env "local" {
  # Database URL for local development
  url = "postgres://postgres:password@localhost:5432/document_reader_chatbot?sslmode=disable"
  
  # Define the development database (used for generating migrations)
  dev = "postgres://postgres:password@localhost:5432/document_reader_chatbot_dev?sslmode=disable"
  
  # Migration directory
  migration {
    dir = "file://migrations"
  }
  
  # Enable format and linting
  format {
    migrate {
      diff = "{{ sql . \" \" }}"
    }
  }
}

env "production" {
  # Production database URL (to be set via environment variable)
  url = env("DATABASE_URL")
  
  # Migration directory
  migration {
    dir = "file://migrations"
  }
  
  # Additional production settings
  format {
    migrate {
      diff = "{{ sql . \" \" }}"
    }
  }
  
  # Disable destructive operations in production
  lint {
    destructive {
      error = true
    }
  }
}

# Schema definition (optional - for GORM integration)
variable "destructive" {
  type    = bool
  default = false
} 