from dataclasses import asdict, is_dataclass
from typing import Any, Type, TypeVar

TOrm = TypeVar("TOrm")
TDomain = TypeVar("TDomain")

def orm_columns(model_cls: Type[Any]) -> set[str]:
    return {c.name for c in model_cls.__table__.columns}  # type: ignore[attr-defined]

def domain_to_orm(domain_obj: TDomain, model_cls: Type[TOrm], *, include_id: bool = False) -> TOrm:
    assert is_dataclass(domain_obj), "domain_to_orm expects a dataclass domain object"
    cols = orm_columns(model_cls)
    data = asdict(domain_obj)
    if not include_id:
        data.pop("id", None)
    payload = {k: v for k, v in data.items() if k in cols}
    return model_cls(**payload)  # type: ignore[misc]

def orm_to_domain(orm_obj: Any, domain_cls: Type[TDomain]) -> TDomain:
    cols = orm_columns(type(orm_obj))
    data = {name: getattr(orm_obj, name) for name in cols}
    return domain_cls(**data)  # type: ignore[misc]

def apply_domain_to_orm(orm_obj: Any, domain_obj: Any, *, skip_id: bool = True) -> None:
    for name, value in asdict(domain_obj).items():
        if skip_id and name == "id":
            continue
        if hasattr(orm_obj, name):
            setattr(orm_obj, name, value)
