import uuid
from app.domain.repositories.book_repo import IBookRepository
from app.usecases.errors import NotFound

class DeleteBookUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self, book_id: uuid.UUID) -> None:
        if self.repo.delete(book_id):
            return None
        else:
            return NotFound("Book not found")