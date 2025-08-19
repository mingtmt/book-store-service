import uuid
from sqlalchemy import select, delete
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session
from app.domain.entities.book import Book
from app.domain.repositories.book_repo import IBookRepository
from app.domain.errors import ConstraintViolation
from app.infrastructure.db.sqlalchemy.models.book_model import BookModel
from app.infrastructure.db.sqlalchemy.mappers.orm_mapper import (
    domain_to_orm, orm_to_domain, apply_domain_to_orm
)

class SqlAlchemyBookRepository(IBookRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_id(self, id: uuid.UUID) -> Book | None:
        stmt = select(BookModel).where(BookModel.id == id)
        m = self.db.execute(stmt).scalars().first()
        return orm_to_domain(m, Book) if m else None

    def get_all(self) -> list[Book] | None:
        rows = self.db.execute(select(BookModel)).scalars().all()
        return [orm_to_domain(r, Book) for r in rows] if rows else None

    def create(self, book: Book) -> Book:
        m = domain_to_orm(book, BookModel)
        self.db.add(m)
        try:
            self.db.commit()
            self.db.refresh(m)
        except IntegrityError as e:
            self.db.rollback()
            # translate DB -> domain error (để middleware map 409)
            raise ConstraintViolation("Book violates a DB constraint") from e
        return orm_to_domain(m, Book)
    
    def save(self, book: Book) -> Book:
        m = self.db.get(BookModel, book.id)
        is_new = m is None
        if is_new:
            m = domain_to_orm(book, BookModel, include_id=False)
            self.db.add(m)
        else:
            apply_domain_to_orm(m, book)

        try:
            self.db.commit()
            self.db.refresh(m)
        except IntegrityError as e:
            self.db.rollback()
            raise ConstraintViolation("Book violates a DB constraint") from e
        return orm_to_domain(m, Book)

    def delete(self, id: uuid.UUID) -> bool:
        result = self.db.execute(delete(BookModel).where(BookModel.id == id))
        self.db.commit()
        return bool(getattr(result, "rowcount", 0))