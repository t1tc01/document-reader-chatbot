from pydantic import BaseModel
from typing import List, Optional

class CreateDocumentModel(BaseModel):
    collection_name: str
    files_url: List[str]
    language: Optional[str] = 'EN'

class DeleteDocumentModel(BaseModel):
    collection_name: str
    document_name: str

