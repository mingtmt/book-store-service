import uuid
from sqlalchemy import Column, String, Index
from sqlalchemy.dialects.postgresql import UUID
from app.core.database import Base

class UserModel(Base):
    __tablename__ = "users"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4, unique=True, nullable=False)
    email = Column(String, unique=True, index=True, nullable=False)
    hashed_password = Column(String, nullable=False)

Index("uq_users_email", UserModel.email, unique=True)