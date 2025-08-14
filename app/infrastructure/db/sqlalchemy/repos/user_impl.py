import uuid
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session
from app.domain.repos.user_repo import IUserRepository
from app.domain.models.user import User
from app.infrastructure.db.sqlalchemy.models.user_model import UserModel
from app.core.exceptions import EmailAlreadyExistsException

class SqlAlchemyUserRepository(IUserRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_email(self, email: str) -> User | None:
        db_user = self.db.query(UserModel).filter(UserModel.email == email).first()
        if db_user:
            return User(id=db_user.id, email=db_user.email, hashed_password=db_user.hashed_password)
        return None

    def create(self, user: User) -> User:
        db_user = UserModel(
            email=user.email,
            hashed_password=user.hashed_password
        )
        self.db.add(db_user)
        try:
            self.db.commit()
        except IntegrityError:
            self.db.rollback()
            raise EmailAlreadyExistsException()
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
