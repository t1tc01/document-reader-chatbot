import os
from model.document import CreateDocumentModel, DeleteDocumentModel
from loguru import logger
from fastapi import HTTPException
from internal.common.preprocessing import process_documents
from config.app_setting import settings
from core.qdrant import q_client
from core.llm import llm_embeddings
from langchain.vectorstores import Qdrant

async def create_documents(documents: CreateDocumentModel):
    """
        Add new document into collection, create new collection if not exists
    """
    collection_name = documents.collection_name
    language = documents.language
    
    files_url = []
    # process file in local 
    for file in documents.files_url:
        if not os.path.isfile(file):
            logger.error(f"File not found: {file}")
            raise HTTPException(402, detail=f"File not found: {file}")
        else:
            files_url.append(file)

    #Begin split documents
    logger.info("Processing documents...")
    try:
        logger.info(f"Split documents...")
        documents_split = process_documents(files_url, [], language)

        ## save preprocessed text to debug
        all_text = " \n ======>>>>>> \n".join([doc.page_content for doc in documents_split])
        save_txt_file = f"{settings.source_documents}/{collection_name}.txt"
        logger.info(f"Save text to {save_txt_file}")
        with open(save_txt_file, "w+") as f:
            f.write(all_text)
    except Exception as e:
        logger.error(f"Split doc error: {e}")
        raise HTTPException(422, detail=str(e))
    
    try:
        ## add documents to vectordb Qdrant 
        logger.info(f"Creating embeddings. May take some minutes...")
        Qdrant.from_documents(
            documents=documents_split,
            embedding=llm_embeddings,
            location=settings.vectordb_qdrant_url,
            collection_name=collection_name
        )
    except Exception as e:
        logger.error(f"Ingest doc failed: {e}")
        raise HTTPException(422, detail=str(e))

    logger.info(f"Ingestion complete!")

async def delete_document(req: DeleteDocumentModel):
    """
        Delete document in collection by source name
    """
    logger.info("Delete document in vector database")
    records = q_client.scroll(collection_name=req.collection_name, limit=1000)

    del_vectors_id = []
    for record in records[0]:
        payload = record.payload
        point_id = record.id
        source = payload['metadata']['source']
        if os.path.basename(req.source_name) == os.path.basename(source):
            del_vectors_id.append(point_id)

    try:
        res = q_client.delete(collection_name=req.collection_name, points_selector=del_vectors_id)
        return {'res': res.status}  
    except Exception as e:
        logger.error(f"Delete document {req.source_name} error in collection {req.collection_name}")
        logger.error(e)
        return {'res': 500}