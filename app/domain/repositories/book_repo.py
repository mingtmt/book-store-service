import uuid
from abc import ABC, abstractmethod

from app.domain.entities.book import Book


class IBookRepository(ABC):
    @abstractmethod
    def get_by_id(self, id: uuid.UUID) -> Book | None:
        pass

    @abstractmethod
    def get_all(self) -> list[Book] | None:
        pass

    @abstractmethod
    def create(self, book: Book) -> Book:
        pass

    @abstractmethod
    def save(self, book: Book) -> Book:
        pass

    @abstractmethod
    def delete(self, id: uuid.UUID) -> bool:
        pass
