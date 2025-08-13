
import uuid
from app.domain.repos.user_repo import UserRepository
from app.core.exceptions import UserNotFoundException

class DeleteUserUseCase:
    def __init__(self, user_repo: UserRepository):
        self.user_repo = user_repo

    def execute(self, user_id: uuid.UUID) -> None:
        if self.user_repo.delete(user_id):
            return None
        else:
            raise UserNotFoundException()