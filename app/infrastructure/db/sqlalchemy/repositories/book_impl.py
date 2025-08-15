import uuid
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session
from app.domain.repositories.book_repo import IBookRepository
from app.domain.entities.book import Book
from app.infrastructure.db.sqlalchemy.models.book_model import BookModel
from fastapi import HTTPException

class SqlAlchemyBookRepository(IBookRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_id(self, id: uuid.UUID) -> Book | None:
        pass

    def get_all(self) -> list[Book]:
        pass

    def create(self, book: Book) -> Book:
        db_book = BookModel(
            title=book.title,
            author=book.author,
            price=book.price,
            description=book.description,
            category=book.category,
        )
        self.db.add(db_book)
        try:
            self.db.commit()
            self.db.refresh(db_book)
        except IntegrityError as e:
            self.db.rollback()
            raise HTTPException(status_code=409, detail="Book already exists") from e

        return Book(
            id=db_book.id,
            title=db_book.title,
            author=db_book.author,
            price=db_book.price,
            description=db_book.description,
            category=db_book.category,
        )

    def update(self, book: Book) -> Book:
        pass

    def delete(self, id: uuid.UUID) -> bool:
        pass