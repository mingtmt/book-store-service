import uuid
from dataclasses import dataclass
from decimal import Decimal
from typing import Optional

from app.domain.entities.book import Book
from app.domain.errors import BookNotFound, ConstraintViolation
from app.domain.repositories.book_repo import IBookRepository


@dataclass(frozen=True)
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
            raise BookNotFound(context={"book_id": str(book_id)})

        title = cmd.title.strip() if cmd.title is not None else book.title
        author = cmd.author.strip() if cmd.author is not None else book.author
        description = (
            cmd.description.strip() if cmd.description is not None else book.description
        )
        category = cmd.category.strip() if cmd.category is not None else book.category
        price = cmd.price if cmd.price is not None else book.price

        if not title or not author:
            raise ConstraintViolation("title/author cannot be empty")
        if price is None or price <= 0:
            raise ConstraintViolation("price must be > 0")

        book.title = title
        book.author = author
        book.description = description
        book.category = category
        book.price = price

        return self.repo.save(book)
