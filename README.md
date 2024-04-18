# GetBlock

## Task description
Develop a service which provide to get top 5 addresses by Activity metric for 100 last blocks.<br>
Metric increases per each transaction between addresses.

#### Example
Alice has sent 100 USDC
Alice received 5 DAI

Alice's Activity metric equal 2

## Setting up

clone repo
```sh
  git clone git@github.com:ur-os/getBockStd.git  
  cd getBockStd
```
for the next step you need to have an installed docker
https://docs.docker.com/engine/install/ubuntu/ (if you using ubuntu)
```sh
    sudo docker build -t get-block .
```
This service have a configs. So next you can config service in ``.env`` file, or set all env's on ``docker run`` <br><br>
With config by file:
```sh
sudo docker run --rm -it -p 127.0.0.1:8080:8080 --env-file .env.example get-block:latest
```
By console (don't forget change ports in `-p` if different):
```sh
sudo docker run --rm -it -p 127.0.0.1:8080:8080 -e GET_BLOCK_ENDPOINT=go.getblock.io -e GET_BLOCK_PORT=8080 get-block:latest
```

Your service ready to serve.<br><br>
Try to get a mEtRiC (for example):
```sh
curl http://localhost:8080/getTopFiveUserActivity

#output: {"status":"ok","code":200,"payload":["0x83c41363cbee0081dab75cb841fa24f3db46627e = 533","0xf89d7b9c864f589bbf53a82105107622b35eaa40 = 124","0x58edf78281334335effa23101bbe3371b6a36a51 = 66","0x28c6c06298d514db089934071355e5743bf21d60 = 62","0x75e89d5979e4f6fba9f97c104c2f0afb3f1dcb88 = 61"]}
```
### Configs 
- `GET_BLOCK_API_KEY` - can be passed. if your rpc node use path param for api key, put it here  (example: `1234567890abcde000000000000000000`)
- `GET_BLOCK_DEPTH` - by default it sets `100`. Indicates how many blocks you will use for metric
- `GET_BLOCK_ENDPOINT` - **necessary** to be set. Hostname of your node (example : `go.getblock.io`)
- `GET_BLOCK_PORT` - by default sets `8080`. Port for your application
- `GET_BLOCK_PULLING_STEP` - by default sets `50`. Degree of parallel processing. For `50` - 100 block in 2 batches  
- `GET_BLOCK_TIMEOUT` - by default `300s`. How long you can await metric response
- `GET_BLOCK_RPS_LIMIT` - by default `50ms`. Delay between 2 threads request node

you can find example of `.env` in `.env.example`