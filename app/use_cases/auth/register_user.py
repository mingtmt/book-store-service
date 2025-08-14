from passlib.context import CryptContext
from app.domain.repos.user_repo import IUserRepository
from app.domain.models.user import User
from app.domain.services.token_service import ITokenService
from app.utils.helper import normalize_email

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class RegisterUserUseCase:
    def __init__(self, user_repo: IUserRepository, token_service: ITokenService, issue_token: bool = True):
        self.user_repo = user_repo
        self.token_service = token_service
        self.issue_token = issue_token

    def execute(self, email: str, password: str) -> tuple[User, str | None]:
        normalized_email = normalize_email(email)
        hashed = pwd_context.hash(password)

        new_user = User(id=None, email=normalized_email, hashed_password=hashed)
        saved = self.user_repo.create(new_user)

        token = self.token_service.create_access_token(str(saved.id)) if self.issue_token else None
        return saved, token
