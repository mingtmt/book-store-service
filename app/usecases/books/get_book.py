import uuid
from app.domain.entities.book import Book
from app.domain.repositories.book_repo import IBookRepository

class GetAllBooksUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self) -> list[Book] | None:
        return self.repo.get_all()
    
class GetBookByIdUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self, book_id: uuid.UUID) -> Book | None:
        return self.repo.get_by_id(book_id)
