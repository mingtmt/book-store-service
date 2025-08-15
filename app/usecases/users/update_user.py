import uuid
from passlib.context import CryptContext
from app.domain.repositories.user_repo import IUserRepository
from app.domain.entities.user import User

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class UpdateUserUseCase:
    def __init__(self, user_repo: IUserRepository):
        self.user_repo = user_repo

    def execute(self, id: uuid.UUID, email: str, password: str) -> User:
        hashed = pwd_context.hash(password) if password else None

        updated_user = User(id=id, email=email, hashed_password=hashed)
        return self.user_repo.update(updated_user)