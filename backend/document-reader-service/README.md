
## To export requirement clearly
```
pip install pipreqs
```
or you can use 
```
pip freeze > requirements.txt 
```

## Run service
```
    cd backend/document-reader-service
    gunicorn app:app -k uvicorn.workers.UvicornWorker --timeout 1500 -c ./gunicorn-config.py --reload
```

## Build dockerfile
```
    docker build --file Dockerfile_document-reader -t document_reader:test ..   
```