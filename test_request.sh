#!/bin/bash

callScore() {
  curl "127.0.0.1:7350/v2/rpc/ProcessPayloadRpc" -H "Authorization: Bearer $NAKAMA_USER_TOKEN" --data '"{\"type\": \"score\", \"version\": \"1.0.0\", \"hash\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}"'
}

callScoreDefaultVersion() {
  curl "127.0.0.1:7350/v2/rpc/ProcessPayloadRpc" -H "Authorization: Bearer $NAKAMA_USER_TOKEN" --data '"{\"type\": \"score\", \"hash\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}"'
}

callScoreDefaultHash() {
  curl "127.0.0.1:7350/v2/rpc/ProcessPayloadRpc" -H "Authorization: Bearer $NAKAMA_USER_TOKEN" --data '"{\"type\": \"score\", \"version\": \"1.0.0\"}"'
}

callScoreMissingFile() {
  curl "127.0.0.1:7350/v2/rpc/ProcessPayloadRpc" -H "Authorization: Bearer $NAKAMA_USER_TOKEN" --data '"{\"type\": \"score\", \"version\": \"2.0.0\", \"hash\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}"'
}

callWithDefaultType() {
  curl "127.0.0.1:7350/v2/rpc/ProcessPayloadRpc" -H "Authorization: Bearer $NAKAMA_USER_TOKEN" --data '"{\"version\": \"1.0.0\", \"hash\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}"'
}

callWithUnknownType() {
  curl "127.0.0.1:7350/v2/rpc/ProcessPayloadRpc" -H "Authorization: Bearer $NAKAMA_USER_TOKEN" --data '"{\"type\": \"friendship\", \"version\": \"1.0.0\", \"hash\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"}"'
}

callWithEmptyPayload() {
  curl "127.0.0.1:7350/v2/rpc/ProcessPayloadRpc" -H "Authorization: Bearer $NAKAMA_USER_TOKEN" --data '"{}"'
}

export NAKAMA_USER_TOKEN=$(curl "127.0.0.1:7350/v2/account/authenticate/device" --data "{\"id\": \""$(uuidgen)"\"}" --user 'defaultkey:' | jq -r '.token')

if [ "$#" -eq 1 ]; then
    echo "Parameter 'callScore' provided. Calling callScore function..."
    echo $("$1")
elif [ "$#" -eq 0 ]; then
    echo "No parameter provided."
else
    echo "Unknown parameter '$1'."
fi