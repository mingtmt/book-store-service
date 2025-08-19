import uuid
from dataclasses import dataclass
from typing import Optional
from decimal import Decimal
from app.domain.entities.book import Book
from app.domain.repositories.book_repo import IBookRepository
from app.domain.errors import BookNotFound, ConstraintViolation

@dataclass
class UpdateBookCommand:
    title: Optional[str] = None
    author: Optional[str] = None
    price: Optional[Decimal] = None
    description: Optional[str] = None
    category: Optional[str] = None

class UpdateBookUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self, book_id: uuid.UUID, cmd: UpdateBookCommand) -> Book:
        book = self.repo.get_by_id(book_id)
        if not book:
            raise BookNotFound()

        if cmd.title is not None:
            book.title = cmd.title
        if cmd.author is not None:
            book.author = cmd.author
        if cmd.price is not None:
            book.price = cmd.price
        if cmd.description is not None:
            book.description = cmd.description
        if cmd.category is not None:
            book.category = cmd.category

        if not book.title or not book.author:
            raise ConstraintViolation("title/author cannot be empty")
        if book.price <= 0:
            raise ConstraintViolation("price must be > 0")

        return self.repo.save(book)
