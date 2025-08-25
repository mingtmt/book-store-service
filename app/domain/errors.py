from __future__ import annotations

from typing import Any, Mapping


class DomainError(Exception):
    code: str = "DOMAIN_ERROR"
    default_message: str = "Domain error"

    def __init__(
        self,
        message: str | None = None,
        *,
        context: Mapping[str, Any] | None = None,
        cause: Exception | None = None,
    ) -> None:
        self.message = message or self.default_message
        self.context = dict(context or {})
        if cause is not None:
            self.__cause__ = cause  # for exception chaining
        super().__init__(self.message)

    def to_dict(self) -> dict[str, Any]:
        data = {"code": self.code, "message": self.message}
        if self.context:
            data["context"] = self.context
        return data

    def __str__(self) -> str:  # nicer string
        return self.message

    def __repr__(self) -> str:
        return f"{self.__class__.__name__}(code={self.code!r}, message={self.message!r}, context={self.context!r})"


# ------- Error families ---------


class NotFoundError(DomainError):
    code = "NOT_FOUND"
    default_message = "Resource not found"


class ConflictError(DomainError):
    code = "CONFLICT"
    default_message = "Conflict"


class ValidationError(DomainError):
    code = "VALIDATION_ERROR"
    default_message = "Validation error"


class AuthError(DomainError):
    code = "AUTH_ERROR"
    default_message = "Authentication/Authorization error"


# ------- Specific domain errors


class UserNotFound(NotFoundError):
    code = "USER_NOT_FOUND"
    default_message = "User not found"


class BookNotFound(NotFoundError):
    code = "BOOK_NOT_FOUND"
    default_message = "Book not found"


class RoleNotFound(NotFoundError):
    code = "ROLE_NOT_FOUND"
    default_message = "Role not found"


class EmailAlreadyExists(ConflictError):
    code = "EMAIL_ALREADY_EXISTS"
    default_message = "Email already exists"


class ConstraintViolation(ConflictError):
    code = "CONSTRAINT_VIOLATION"
    default_message = "Resource violates data constraints"


class InvalidCredentials(AuthError):
    code = "INVALID_CREDENTIALS"
    default_message = "Invalid email or password"


class AccessTokenMissing(AuthError):
    code = "ACCESS_TOKEN_ERROR"
    default_message = "Access token is missing"


class AccessTokenInvalid(AuthError):
    code = "ACCESS_TOKEN_ERROR"
    default_message = "Access token is invalid"


class AccessTokenExpired(AuthError):
    code = "ACCESS_TOKEN_ERROR"
    default_message = "Access token is expired"


class InvalidEmail(ValidationError):
    code = "INVALID_EMAIL"
    default_message = "Email format is invalid"


class PasswordTooWeak(ValidationError):
    code = "PASSWORD_TOO_WEAK"
    default_message = "Password does not meet strength requirements"
