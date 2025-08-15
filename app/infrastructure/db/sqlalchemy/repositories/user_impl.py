import uuid
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session
from sqlalchemy import func
from app.domain.repositories.user_repo import IUserRepository
from app.domain.entities.user import User
from app.domain.errors import EmailAlreadyExists
from app.infrastructure.db.sqlalchemy.models.user_model import UserModel
from app.utils.helper import normalize_email

class SqlAlchemyUserRepository(IUserRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_email(self, email: str) -> User | None:
        normalized = normalize_email(email)
        db_user = self.db.query(UserModel).filter(func.lower(UserModel.email) == normalized).first()
        return User(id=db_user.id, email=db_user.email, hashed_password=db_user.hashed_password) if db_user else None

    def create(self, user: User) -> User:
        db_user = UserModel(email=user.email, hashed_password=user.hashed_password)
        self.db.add(db_user)
        try:
            self.db.commit()
        except IntegrityError as e:
            self.db.rollback()
            raise EmailAlreadyExists(str(e)) from e
        self.db.refresh(db_user)
        return User(id=db_user.id, email=db_user.email, hashed_password=db_user.hashed_password)
    
    def update(self, user: User) -> User:
        db_user = self.db.query(UserModel).filter(UserModel.id == user.id).first()
        if not db_user:
            return None
        db_user.email = user.email if user.email else db_user.email
        if user.hashed_password:
            db_user.hashed_password = user.hashed_password
        self.db.commit()
        return User(id=db_user.id, email=db_user.email, hashed_password=db_user.hashed_password)

    def delete(self, id: uuid.UUID) -> bool:
        db_user = self.db.query(UserModel).filter(UserModel.id == id).first()
        if db_user:
            self.db.delete(db_user)
            self.db.commit()
            return True
        return False
