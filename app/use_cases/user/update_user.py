import uuid
from passlib.context import CryptContext
from app.domain.repos.user_repo import UserRepository
from app.domain.models.user import User

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class UpdateUserUseCase:
    def __init__(self, user_repo: UserRepository):
        self.user_repo = user_repo

    def execute(self, id: uuid.UUID, email: str, password: str) -> User:
        hashed = pwd_context.hash(password) if password else None

        updated_user = User(id=id, email=email, hashed_password=hashed)
        return self.user_repo.update(updated_user)