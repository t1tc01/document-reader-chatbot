
## To export requirement clearly
```
pip install pipreqs
```

## Run service
```
    cd backend/document-reader-service
    gunicorn app:app -k uvicorn.workers.UvicornWorker --timeout 1500 -c ./gunicorn-config.py --reload
```