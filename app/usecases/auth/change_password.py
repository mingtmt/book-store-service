import uuid
from dataclasses import dataclass

from app.domain.errors import InvalidCredentials, UserNotFound
from app.domain.repositories.user_repo import IUserRepository
from app.ports.password_service import IPasswordService


@dataclass
class ChangePasswordCommand:
    id: uuid.UUID
    old_password: str
    new_password: str


class ChangePasswordUseCase:
    def __init__(
        self, repo: IUserRepository, password_service: IPasswordService
    ) -> None:
        self.repo = repo
        self.password_service = password_service

    def execute(self, user_id: uuid.UUID, cmd: ChangePasswordCommand) -> None:
        user = self.repo.get_by_id(user_id)
        if not user:
            raise UserNotFound(context={"user_id": user_id})

        if not self.password_service.verify_password(
            cmd.old_password, user.hashed_password
        ):
            raise InvalidCredentials("Old password is incorrect")
        user.hashed_password = self.password_service.hash_password(cmd.new_password)
        self.repo.save(user)
