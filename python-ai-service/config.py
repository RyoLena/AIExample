import os
from dotenv import load_dotenv
class Settings:
    load_dotenv()
    model_name = os.environ.get("MODEL_NAME")
    port = int(os.environ.get("PORT"))
    gemin_api_key = os.environ.get("GOOGLE_API_KEY")
setting = Settings()

