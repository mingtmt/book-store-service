import uuid

from sqlalchemy import func, select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session

from app.domain.entities.user import User
from app.domain.errors import ConstraintViolation, EmailAlreadyExists
from app.domain.repositories.user_repo import IUserRepository
from app.infrastructure.db.sqlalchemy.mappers.orm_mapper import (
    apply_domain_to_orm,
    domain_to_orm,
    orm_to_domain,
)
from app.infrastructure.db.sqlalchemy.models.user_model import UserModel
from app.utils.helper import normalize_email


class SqlAlchemyUserRepository(IUserRepository):
    def __init__(self, db: Session):
        self.db = db

    def create(self, user: User) -> User:
        m = domain_to_orm(user, UserModel)
        self.db.add(m)
        try:
            self.db.commit()
            self.db.refresh(m)
        except IntegrityError as e:
            self.db.rollback()
            msg = str(getattr(e, "orig", e))
            if "uq_users_email_ci" in msg:
                raise ConstraintViolation("Email already exists", cause=e)
            raise ConstraintViolation("Resource violates data constraints", cause=e)
        return orm_to_domain(m, User)

    def get_by_id(self, id: uuid.UUID) -> User | None:
        pass

    def get_by_email(self, email: str) -> User | None:
        normalized_email = normalize_email(email)
        m = (
            self.db.execute(
                select(UserModel).where(
                    func.lower(UserModel.email) == func.lower(normalized_email)
                )
            )
            .scalars()
            .first()
        )
        return orm_to_domain(m, User) if m else None

    def update(self, user: User) -> User:
        db_user = self.db.query(UserModel).filter(UserModel.id == user.id).first()
        if not db_user:
            return None
        db_user.email = user.email if user.email else db_user.email
        if user.hashed_password:
            db_user.hashed_password = user.hashed_password
        self.db.commit()
        return User(
            id=db_user.id, email=db_user.email, hashed_password=db_user.hashed_password
        )

    def delete(self, id: uuid.UUID) -> bool:
        db_user = self.db.query(UserModel).filter(UserModel.id == id).first()
        if db_user:
            self.db.delete(db_user)
            self.db.commit()
            return True
        return False
