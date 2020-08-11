from fastapi import FastAPI, Body
from pydantic import BaseModel

class Value(BaseModel):
    value: int

app = FastAPI()

@app.post('/')
def root(v: Value):
    print('value:', v.value)
    return {'response': v.value**2}
