
import uuid
from app.domain.repositories.user_repo import IUserRepository
from app.usecases.errors import NotFound

class DeleteUserUseCase:
    def __init__(self, user_repo: IUserRepository):
        self.user_repo = user_repo

    def execute(self, user_id: uuid.UUID) -> None:
        if self.user_repo.delete(user_id):
            return None
        else:
            raise NotFound("User not found")