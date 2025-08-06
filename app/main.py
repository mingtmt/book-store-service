from fastapi import FastAPI
from app.api.v1.routes import ping

app = FastAPI()

# Mount route
app.include_router(ping.router, prefix="/api/v1", tags=["Ping"])

