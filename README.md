# deviceshifu-ai-sfn

1. Install YoMo cli

```sh
go install github.com/yomorun/yomo/cmd/yomo@latest
```

2. Start mock server

```sh
uvicorn mock_server:app --port 30080
```

3. Get vivgrid.com token

# https://dashboard.vivgrid.com

```sh
export VIVGRID_TOKEN="********"
```

4. Run get_image SFN

```sh
cd get_image

export YOMO_SFN_CREDENTIAL="app-key-secret:$VIVGRID_TOKEN"

yomo run app.go
```

```sh
curl https://openai.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $VIVGRID_TOKEN" \
  -d '{
     "model": "gpt-4o-mini",
     "messages": [{"role": "user", "content": "Hi, can you tell what the camera is seeing?"}]
   }'
```

```sh
open ./image.png
```
