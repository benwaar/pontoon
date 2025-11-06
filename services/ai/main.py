from fastapi import FastAPI

app = FastAPI()


@app.get("/healthz")
def healthz():
    return {"status": "ok"}


@app.get("/ai/hint")
def hint():
    return {"hint": "Consider sticking at 18."}
