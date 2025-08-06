from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from app.api.dependencies.db import get_db
from app.domain.models.ping import Ping

router = APIRouter()

@router.get("/ping")
def ping(db: Session = Depends(get_db)):
    ping = db.query(Ping).first()
    return {"message": ping.message if ping else "no data"}
