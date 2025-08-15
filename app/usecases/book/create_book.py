from decimal import Decimal
from app.domain.repositories.book_repo import IBookRepository
from app.domain.entities.book import Book

class CreateBookUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self, title: str, author: str, price: Decimal, description: str, category: str) -> Book:
        new_book = Book(id=None, title=title, author=author, price=price, description=description, category=category)
        saved = self.repo.create(new_book)
        return saved
