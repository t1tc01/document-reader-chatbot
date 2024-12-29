from langchain.docstore.document import Document
from langchain.document_loaders import (
    PyMuPDFLoader,
)
from typing import List
from loguru import logger
import os
from multiprocessing import Pool
import tqdm

LOADER_MAPPING = {
    ".pdf": (PyMuPDFLoader, {}),
}

#
def load_documents(process_files: List[str], ignored_files: List[str] = []) -> List[Document]:
    """
    Loads all documents, ignoring specified files
    """
    with Pool(processes=os.cpu_count()) as pool:
        results = []
        with tqdm(total=len(process_files), desc='Loading new documents', ncols=80) as pbar:
            for i, doc in enumerate(pool.imap_unordered(load_single_document, process_files)):
                if doc is not None:
                    results.append(doc)
                pbar.update()
    return results


def load_single_document(file_path: str) -> List[Document]:
    """
        Load single document from file path
    """
    ext = "." + file_path.rsplit(".", 1)[-1]

    if ext in LOADER_MAPPING:
        logger.info(f"Loading file: {file_path}")
        loader_class, loader_args = LOADER_MAPPING[ext]
        loader = loader_class(file_path, **loader_args)
        docs = loader.load()

        try:
            doc = docs[0]
            doc.page_content = get_all_text_in_doc(docs)
            doc.metadata["source"] = file_path
            doc.metadata["name"] = os.path.basename(file_path)
            return doc
        except Exception as e:
            logger.error("Error in load single doc")
            return None
    raise ValueError(f"Unsupported file extension '{ext}'")

def get_all_text_in_doc(docs):
    all_text = [doc.page_content for doc in docs]
    all_text = " \n".join(all_text)
    return all_text

def process_documents(
        process_files: List[str],
        ignored_files: List[str] = [],
        language='EN'
    ) -> List[Document]:
    """
    Process documents
    """
    logger.info("Processing documents...")
    if language == "EN":
        docs = process_documents_en(process_files, ignored_files)
    return docs

def process_documents_en(process_files: List[str], ignored_files: List[str] = []) -> List[Document]:
    logger.info("Process documents English")