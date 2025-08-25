import uuid
from typing import List

from sqlalchemy import Index, Integer, String, func
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.infrastructure.db.session import Base
from app.infrastructure.db.sqlalchemy.models.role_model import RoleModel, user_roles


class UserModel(Base):
    __tablename__ = "users"

    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid.uuid4,
        nullable=False,
    )

    email: Mapped[str] = mapped_column(String(320), nullable=False)
    name: Mapped[str] = mapped_column(String(100), nullable=False)
    age: Mapped[int] = mapped_column(Integer, nullable=False)
    hashed_password: Mapped[str] = mapped_column(String(255), nullable=False)

    roles: Mapped[List[RoleModel]] = relationship(
        "RoleModel",
        secondary=user_roles,
        backref="users",
        lazy="selectin",
        cascade="save-update",
    )

    __table_args__ = (Index("uq_users_email_ci", func.lower(email), unique=True),)

    def __repr__(self):
        return f"<User(id={self.id}, email={self.email})>"
