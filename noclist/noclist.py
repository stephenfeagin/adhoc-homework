#!python3

import hashlib
import json
import urllib3

import requests


def retry(url, max_tries=3, **kwargs):
    n_tries = 0
    while n_tries < max_tries:
        n_tries += 1
        try:
            response = requests.get(url, **kwargs)
        except requests.exceptions.ConnectionError:
            continue
        if response.status_code == requests.codes.ok:
            return response
    raise requests.exceptions.HTTPError


def get_auth_token(url):
    try:
        response = retry(url + "/auth")
    except requests.exceptions.HTTPError:
        return 1
    token = response.headers["Badsec-Authentication-Token"]
    return token


def get_auth_checksum(token):
    encoded = (token + "/users").encode("utf-8")
    checksum = hashlib.sha256(encoded).hexdigest()
    return checksum


def get_user_ids(url, checksum):
    response = retry(url + "/users", headers={"X-Request-Checksum": checksum})
    return response.content.decode()


def get_badsec_users(url="http://127.0.0.1:8888"):
    token = get_auth_token(url=url)
    checksum = get_auth_checksum(token=token)
    users = get_user_ids(url=url, checksum=checksum)
    print(json.dumps(users.splitlines()))
    exit(0)


if __name__ == '__main__':
    get_badsec_users()
