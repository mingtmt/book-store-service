def normalize_email(email: str) -> str:
    import unicodedata

    return unicodedata.normalize("NFC", email.strip()).casefold()
