from fastapi import FastAPI
from fastapi.responses import FileResponse

app = FastAPI()


@app.get('/deviceshifu-camera/capture')
async def chat_completions():
    return FileResponse('./mock.png')
