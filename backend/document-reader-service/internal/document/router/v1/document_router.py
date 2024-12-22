from fastapi import APIRouter
from model.document import CreateDocumentModel, DeleteDocumentModel
from loguru import logger

router = APIRouter()

@router.post("/add")
async def add_documents(documents: CreateDocumentModel):
    logger.info(f'Call API add document')
    return {"message": "Files embedded successfully"}

@router.delete("/delete")
async def delete_document():
    logger.info(f'Call API delete document')
    return {'message':'File delete successfully'}