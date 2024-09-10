from random import randint
from fastapi import FastAPI
from fastapi.responses import FileResponse

app = FastAPI()


@app.get('/deviceshifu-camera/capture')
async def camera():
    return FileResponse('./mock.jpg')


@app.get('/deviceshifu-plc/sendsinglebit')
async def plc(value: int):
    return {'switch': 'OPEN' if value else 'CLOSED'}


@app.get('/deviceshifu-led/number')
async def led():
    return {'number': randint(0, 9)}
