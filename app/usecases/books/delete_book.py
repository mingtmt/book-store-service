import uuid
from app.domain.repositories.book_repo import IBookRepository
from app.domain.errors import BookNotFound

class DeleteBookUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self, book_id: uuid.UUID) -> None:
        deleted = self.repo.delete(book_id)
        if not deleted:
            raise BookNotFound(context={"book_id": str(book_id)})
        else:
            return None