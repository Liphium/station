### Note

LiveKit lost support in the last release. We're looking to rework the entire implementation of Spaces and getting rid of LiveKit is just the first step for the desire for more control. Expect voice chat to return in future versions of Liphium, but that will likely take months to years until arrival.

What this means for you and why this notice is here is that you can now safely remove all environment variables in the previous environment file starting with the SS_LK prefix. They are no longer needed.

### Example .env file for this app

```bash
# BACKEND CONFIGURATION

# Domain config
BASE_PATH = "localhost:3000"
BASE_PORT = "3000"
CHAT_NODE = "localhost:3001"
CHAT_NODE_PORT = "3001"
SPACE_NODE = "localhost:3002"
SPACE_NODE_PORT = "3002"

# App config
APP_NAME = "Liphium"
TESTING = "true"
LISTEN = "127.0.0.1"
# PROTOCOL = "http://" (You CAN do this, but https is the default)
CLI = "true"
SYSTEM_UUID = "your-uuid" # You can delete this and the app will give you a random one

# Database
DB_USER = "postgres"
DB_PASSWORD = "deinemutter123"
DB_DATABASE = "backend"
DB_HOST = "localhost"
DB_PORT = "5432"

# SSO (if you want it)
# SSO_ENABLED = "true"
# SSO_CONFIG = "url"
# SSO_CLIENT_ID = "id"
# SSO_CLIENT_SECRET = "secret"

# JWT
JWT_SECRET = "secret"

# Through Cloudflare Protection
TC_PUBLIC_KEY = "PUB_KEY_HERE"
TC_PRIVATE_KEY = "PRIV_KEY_HERE"

# File storage folder
# If you want local storage
FILE_REPO_TYPE = "local"
FILE_REPO = "/home/some_path"

# If you want to use R2
# FILE_REPO_TYPE = "r2"
# FILE_REPO_KEY_ID = "id"
# FILE_REPO_KEY = "key"
# FILE_REPO = "https://account.eu.r2.cloudflarestorage.com"
# FILE_REPO_BUCKET = "your-bucket"

# SMTP (for emails)
SMTP_SERVER = "mail.example.com"
SMTP_PORT = "port"
SMTP_FROM = "no-reply@example.com"
SMTP_USER = "username"
SMTP_PW = "password"


# CHAT NODE CONFIGURATION

# Config (If you want to allow unsafe locations (http), YOU SHOULD ONLY FOR TESTING)
# CN_ALLOW_UNSAFE = "true"

# Database
CN_DB_USER = "postgres"
CN_DB_PASSWORD = "deinemutter123"
CN_DB_DATABASE = "chat"
CN_DB_HOST = "localhost"
CN_DB_PORT = "5432"
```
