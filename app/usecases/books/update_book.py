import uuid
from app.domain.entities.book import Book
from app.domain.repositories.book_repo import IBookRepository
from app.presentation.http.schemas.books import UpdateBookRequest

class UpdateBookUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self, book_id: uuid.UUID, payload: UpdateBookRequest) -> Book:
        return self.repo.update(book_id, payload)