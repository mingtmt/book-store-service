import uuid
from decimal import Decimal
from typing import Dict, Optional

import pytest

from app.domain.entities.book import Book
from app.domain.errors import BookNotFound, ConstraintViolation


def _norm(s: Optional[str]) -> Optional[str]:
    if s is None:
        return None
    return " ".join(s.split()).lower()


class InMemoryBookRepo:
    """
    Fake IBookRepository
    Fields: id, title, author, price, description, category
    Constraints:
      - Unique (title, author) case-insensitive
      - price >= 0
    """

    def __init__(self):
        self._data: Dict[uuid.UUID, Book] = {}

    # --- helpers ---
    def _ensure_price_ok(self, price: Decimal):
        if price is None:
            return
        if Decimal(price) < 0:
            raise ConstraintViolation("Price must be non-negative")

    def _assert_unique_title_author(
        self, title: str, author: str, exclude_id: Optional[uuid.UUID] = None
    ):
        t = _norm(title)
        a = _norm(author)
        for bid, b in self._data.items():
            if exclude_id and bid == exclude_id:
                continue
            if _norm(b.title) == t and _norm(b.author) == a:
                raise ConstraintViolation("Title & author must be unique")

    # --- interface methods ---
    def get_by_id(self, id: uuid.UUID) -> Optional[Book]:
        return self._data.get(id)

    def get_all(self) -> list[Book]:
        return list(self._data.values())

    def create(self, book: Book) -> Book:
        self._assert_unique_title_author(book.title, book.author)
        self._ensure_price_ok(book.price)

        new_id = getattr(book, "id", None) or uuid.uuid4()
        created = Book(
            id=new_id,
            title=book.title,
            author=book.author,
            price=book.price,
            description=book.description,
            category=book.category,
        )
        self._data[new_id] = created
        return created

    def save(self, book: Book) -> Book:
        if book.id not in self._data:
            raise BookNotFound(context={"book_id": book.id})
        self._assert_unique_title_author(book.title, book.author, exclude_id=book.id)
        self._ensure_price_ok(book.price)
        updated = Book(
            id=book.id,
            title=book.title,
            author=book.author,
            price=book.price,
            description=book.description,
            category=book.category,
        )
        self._data[book.id] = updated
        return updated

    def delete(self, id: uuid.UUID) -> bool:
        if id not in self._data:
            return False
        del self._data[id]
        return True


@pytest.fixture
def fake_book_repo():
    return InMemoryBookRepo()


@pytest.fixture
def book_payload():
    return {
        "title": "Clean Architecture in Action",
        "author": "Robert C. Martin Jr.",
        "price": Decimal("19.99"),
        "description": "A hands-on guide",
        "category": "software",
    }


@pytest.fixture
def user_payload():
    """A default user payload suitable for repo.create(User(**...))."""
    return {
        "email": "ada.lovelace@example.com",
        "name": "Ada Lovelace",
        "age": 28,
        "hashed_password": "HASHED",
    }
