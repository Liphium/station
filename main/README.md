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
TESTING_AMOUNT = "2"
PROTOCOL = "http://"
CLI = "true"

# Database (for the backend)
DB_USER = "postgres"
DB_PASSWORD = "deinemutter123"
DB_DATABASE = "node_backend"
DB_HOST = "localhost"
DB_PORT = "5432"

# JWT
JWT_SECRET = "secret"

# Through Cloudflare Protection (or reverse proxy protection, doesn't help against man in the middle, only trusted reverse proxies)
TC_PUBLIC_KEY="KEY_HERE"
TC_PRIVATE_KEY="KEY_HERE"

# File storage folder (choose some path that exists)
FILE_REPO = "/home/julian/Documents/repo/cloud"

# File upload settings (defaults are fine, make sure to change nginx if modified)
MAX_UPLOAD_SIZE = "10" # in MB
MAX_FAVORITE_STORAGE = "500" # in MB
MAX_TOTAL_STORAGE = "1000" # in MB

# SMTP (for emails)
SMTP_SERVER = "email-smtp.eu-north-1.amazonaws.com"
SMTP_PORT = "2587"
SMTP_IDENTITY = "backend"
SMTP_FROM = "no-reply@liphium.app"
SMTP_USER = "user"
SMTP_PW = "password"


# CHAT NODE CONFIGURATION

# Database (for the chat node)
CN_DB_USER = "postgres"
CN_DB_PASSWORD = "deinemutter123"
CN_DB_DATABASE = "chat"
CN_DB_HOST = "localhost"
CN_DB_PORT = "5432"

# Live share (choose some path that exists)
CN_LS_REPO = "/home/julian/Documents/repo/ls"


# SPACE STATION CONFIGURATION

# LiveKit (should work if you install a local server according to https://docs.livekit.io/realtime/self-hosting/local/)
SS_LK_URL = "http://localhost:7880"
SS_LK_KEY = "devkey"
SS_LK_SECRET = "secret"
```