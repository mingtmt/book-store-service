from sqlalchemy.orm import Session
from app.domain.repos.user_repo import UserRepository
from app.domain.models.user import User
from app.infrastructure.db.sqlalchemy.models.user_model import UserModel  # SQLAlchemy table model

class SqlAlchemyUserRepository(UserRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_email(self, email: str) -> User | None:
        db_user = self.db.query(UserModel).filter(UserModel.email == email).first()
        if db_user:
            return User(id=db_user.id, email=db_user.email, hashed_password=db_user.hashed_password)
        return None

    def create(self, user: User) -> User:
        db_user = UserModel(email=user.email, hashed_password=user.hashed_password)
        self.db.add(db_user)
        self.db.commit()
        self.db.refresh(db_user)
        return User(id=db_user.id, email=db_user.email, hashed_password=db_user.hashed_password)
