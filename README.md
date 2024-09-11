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

```sh
export VIVGRID_TOKEN="********" # get from https://dashboard.vivgrid.com
```

4. Run get_image SFN

```sh
cd get_image

export YOMO_SFN_CREDENTIAL="app-key-secret:$VIVGRID_TOKEN"
export OPENAI_API_KEY="sk-******"
export OPENAI_BASE_URL="https://api.openai.com/v1"

yomo run app.go
```

```sh
curl https://openai.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $VIVGRID_TOKEN" \
  -d '{
     "messages": [{"role": "user", "content": "Hi, can you tell what the camera is seeing?"}]
   }'
```

5. Run set_plc_output SFN

```sh
cd set_plc_output

export YOMO_SFN_CREDENTIAL="app-key-secret:$VIVGRID_TOKEN"

yomo run app.go
```

```sh
curl https://openai.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $VIVGRID_TOKEN" \
  -d '{
     "messages": [{"role": "user", "content": "Can you set the PLC output to true?"}]
   }'
```

6. Run set_led SFN

```sh
cd set_led

export YOMO_SFN_CREDENTIAL="app-key-secret:$VIVGRID_TOKEN"

yomo run app.go
```

```sh
curl https://openai.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $VIVGRID_TOKEN" \
  -d '{
     "messages": [{"role": "user", "content": "Can you set the display number on the LED to 4005?"}]
   }'
```
