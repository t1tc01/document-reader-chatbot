from fastapi import APIRouter
from model.document import CreateDocumentModel, DeleteDocumentModel
from loguru import logger
from internal.document.service.document_service import create_documents, delete_document

router = APIRouter()

@router.post("/add")
async def add_documents(documents: CreateDocumentModel):
    logger.info(f'Call API add document')
    await create_documents(documents)
    return {"message": "Files embedded successfully"}

@router.delete("/delete")
async def delete_document(req: DeleteDocumentModel):
    logger.info(f'Call API delete document')
    await delete_document(req)
    return {'message':'File delete successfully'}