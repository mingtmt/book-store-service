from sqlalchemy import Column, Integer, String
from app.core.database import Base

class Ping(Base):
    __tablename__ = "ping"
    id = Column(Integer, primary_key=True, index=True)
    message = Column(String, default="pong")
