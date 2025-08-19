import uuid
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session
from app.domain.repositories.book_repo import IBookRepository
from app.domain.entities.book import Book
from app.infrastructure.db.sqlalchemy.models.book_model import BookModel
from app.presentation.http.schemas.books import UpdateBookRequest
from fastapi import HTTPException

class SqlAlchemyBookRepository(IBookRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_id(self, id: uuid.UUID) -> Book | None:
        db_book = self.db.query(BookModel).filter(BookModel.id == id).first()
        if db_book:
            return Book(
                id=db_book.id,
                title=db_book.title,
                author=db_book.author,
                price=db_book.price,
                description=db_book.description,
                category=db_book.category,
            )
        return None

    def get_all(self) -> list[Book] | None:
        db_books = self.db.query(BookModel).all()
        if not db_books:
            return None
        return [
            Book(
                id=db_book.id,
                title=db_book.title,
                author=db_book.author,
                price=db_book.price,
                description=db_book.description,
                category=db_book.category,
            )
            for db_book in db_books
        ]

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

    def update(self, book_id: uuid.UUID, payload: UpdateBookRequest) -> Book:
        db_book = self.db.query(BookModel).filter(BookModel.id == book_id).first()
        if not db_book:
            raise HTTPException(status_code=404, detail="Book not found")

        update_data = payload.model_dump(exclude_unset=True)

        if not update_data:
            raise HTTPException(status_code=400, detail="No fields provided to update")

        update_data.pop("id", None)

        for key, value in update_data.items():
            setattr(db_book, key, value)

        try:
            self.db.commit()
            self.db.refresh(db_book)
        except IntegrityError as e:
            self.db.rollback()
            raise HTTPException(status_code=400, detail="Invalid data or constraint violation") from e

        return Book(
            id=db_book.id,
            title=db_book.title,
            author=db_book.author,
            price=db_book.price,
            description=db_book.description,
            category=db_book.category,
        )

    def delete(self, id: uuid.UUID) -> bool:
        db_book = self.db.query(BookModel).filter(BookModel.id == id).first()
        if db_book:
            self.db.delete(db_book)
            self.db.commit()
            return True
        return False