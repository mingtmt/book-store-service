from pydantic import BaseModel
from typing import Generic, TypeVar, Optional, Type, Iterable

T = TypeVar("T")
class PageMeta(BaseModel):
    page: int
    size: int
    total: int
class Envelope(BaseModel, Generic[T]):
    data: Optional[T] = None
    meta: Optional[PageMeta] = None
    message: Optional[str] = None

def to_list(out_cls: Type[T], items: Iterable[object]) -> list[T]:
    return [out_cls.model_validate(it) for it in items]