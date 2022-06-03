# Tweety

Tweets fake news classifier

To run the service localy execute:

```
make compose-up
```

Example classify call:

```
GET localhost:8080/v1/tweets/classify?pageId=111
```

Run tests:

```
make test
```

