from dataclasses import dataclass
from decimal import Decimal
from typing import Optional

from app.domain.entities.book import Book
from app.domain.repositories.book_repo import IBookRepository
from app.usecases.errors import BadRequest


@dataclass(frozen=True)
class CreateBookCommand:
    title: str
    author: str
    price: Decimal
    description: Optional[str]
    category: str


class CreateBookUseCase:
    def __init__(self, repo: IBookRepository):
        self.repo = repo

    def execute(self, cmd: CreateBookCommand) -> Book:
        title = cmd.title.strip()
        author = cmd.author.strip()
        if not title:
            raise BadRequest("Title must not be empty")
        if not author:
            raise BadRequest("Author must not be empty")
        if cmd.price < 0:
            raise BadRequest("Price must be non-negative")

        new_book = Book(
            id=None,
            title=title,
            author=author,
            price=cmd.price,
            description=(cmd.description.strip() if cmd.description else None),
            category=cmd.category.strip(),
        )
        return self.repo.create(new_book)
