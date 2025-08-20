import uuid
from sqlalchemy import select, delete
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session
from app.domain.entities.book import Book
from app.domain.repositories.book_repo import IBookRepository
from app.domain.errors import BookNotFound, ConstraintViolation
from app.infrastructure.db.sqlalchemy.models.book_model import BookModel
from app.infrastructure.db.sqlalchemy.mappers.orm_mapper import (
    domain_to_orm, orm_to_domain, apply_domain_to_orm
)

class SqlAlchemyBookRepository(IBookRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_id(self, id: uuid.UUID) -> Book | None:
        m = self.db.get(BookModel, id)
        return orm_to_domain(m, Book) if m else None

    def get_all(self) -> list[Book]:
        rows = self.db.execute(select(BookModel)).scalars().all()
        return [orm_to_domain(r, Book) for r in rows]

    def create(self, book: Book) -> Book:
        m = domain_to_orm(book, BookModel)
        self.db.add(m)
        try:
            self.db.commit()
            self.db.refresh(m)
        except IntegrityError as e:
            self.db.rollback()
            msg = str(getattr(e, "orig", e))
            if "uq_books_title_author_ci" in msg or "uq_books_title_author" in msg:
                raise ConstraintViolation("Title & author must be unique", cause=e)
            if "ck_books_price_nonnegative" in msg:
                raise ConstraintViolation("Price must be non-negative", cause=e)
            raise ConstraintViolation("Resource violates data constraints", cause=e)
        return orm_to_domain(m, Book)
    
    def save(self, book: Book) -> Book:
        m = self.db.get(BookModel, book.id)
        if m is None:
            raise BookNotFound(context={"book_id": str(book.id)})

        apply_domain_to_orm(m, book)
        try:
            self.db.commit()
            self.db.refresh(m)
        except IntegrityError as e:
            self.db.rollback()
            raise ConstraintViolation("DB constraint violated", cause=e)
        return orm_to_domain(m, Book)

    def delete(self, id: uuid.UUID) -> bool:
        result = self.db.execute(delete(BookModel).where(BookModel.id == id))
        self.db.commit()
        return bool(getattr(result, "rowcount", 0))