from fastapi import FastAPI
from app.api.v1.routes import ping
from app.core.database import Base, engine

app = FastAPI()

# Mount route
app.include_router(ping.router, prefix="/api/v1", tags=["Ping"])

Base.metadata.create_all(bind=engine)

