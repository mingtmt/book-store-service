from __future__ import annotations
import re
import uuid
from decimal import Decimal
from typing import Optional, Annotated
from enum import Enum
from pydantic import BaseModel, Field, ConfigDict, field_validator, model_validator


class BookCategory(str, Enum):
    study = "study"
    novel = "novel"
    scifi = "sci-fi"
    kids = "kids"


Price = Annotated[Decimal, Field(gt=0, max_digits=10, decimal_places=2)]


class CreateBook(BaseModel):
    model_config = ConfigDict(extra="forbid", json_schema_extra={
        "example": {
            "title": "Clean Architecture",
            "author": "Robert C. Martin",
            "price": "19.99",
            "description": "Principles of software architecture.",
            "category": "study",
        }
    })

    title: Annotated[str, Field(min_length=1, max_length=255)]
    author: Annotated[str, Field(min_length=1, max_length=255)]
    price: Price
    description: Optional[Annotated[str, Field(max_length=10_000)]] = None
    category: Annotated[str, Field(min_length=1, max_length=100)]

    @field_validator("title", "author", "category", "description", mode="before")
    @classmethod
    def _strip(cls, v):
        if isinstance(v, str):
            v = v.strip()
            v = re.sub(r"\s+", " ", v)
        return v


class UpdateBook(BaseModel):
    model_config = ConfigDict(extra="forbid")

    title: Optional[Annotated[str, Field(min_length=1, max_length=255)]] = None
    author: Optional[Annotated[str, Field(min_length=1, max_length=255)]] = None
    price: Optional[Annotated[Decimal, Field(gt=0, max_digits=10, decimal_places=2)]] = None
    description: Optional[Annotated[str, Field(max_length=10_000)]] = None
    category: Optional[Annotated[str, Field(min_length=1, max_length=100)]] = None

    @field_validator("title", "author", "category", "description", mode="before")
    @classmethod
    def _strip(cls, v):
        if isinstance(v, str):
            v = v.strip()
        return v

    @model_validator(mode="after")
    def _at_least_one_field(self):
        if all(getattr(self, f) is None for f in ("title", "author", "price", "description", "category")):
            raise ValueError("At least one field must be provided")
        return self


class BookOut(BaseModel):
    model_config = ConfigDict(from_attributes=True, extra="ignore")

    id: uuid.UUID
    title: str
    author: str
    price: Decimal
    description: Optional[str] = None
    category: str

    @classmethod
    def from_domain(cls, book) -> "BookOut":
        return cls.model_validate(book)
