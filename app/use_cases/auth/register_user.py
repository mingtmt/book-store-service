from passlib.context import CryptContext
from app.domain.repos.user_repo import UserRepository
from app.domain.models.user import User
from app.infrastructure.services.jwt_service import create_access_token

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class RegisterUserUseCase:
    def __init__(self, user_repo: UserRepository, issue_token: bool = True):
        self.user_repo = user_repo
        self.issue_token = issue_token

    def execute(self, email: str, password: str) -> tuple[User, str | None]:
        exist = self.user_repo.get_by_email(email)
        if exist:
            raise ValueError("EMAIL_ALREADY_EXISTS")

        hashed = pwd_context.hash(password)

        new_user = User(id=None, email=email, hashed_password=hashed)
        saved = self.user_repo.create(new_user)

        token = create_access_token({"user_id": str(saved.id)}) if self.issue_token else None
        return saved, token
