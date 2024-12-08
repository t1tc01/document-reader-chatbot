from config.app_setting import settings

bind = "0.0.0.0:8001"
workers = settings.worker_nb #multiprocessing.cpu_count() * 2 + 1
