#!/usr/bin/env python3
# -*- coding:utf-8 -*-


import hashlib
import hmac
import base64
import json
from urllib import parse
import requests
import time
import uuid


ecs_url = 'https://ecs.aliyuncs.com/'
ram_url = 'https://ram.aliyuncs.com'
access_key = "access_key"
access_key_secret = "access_key_secret"
FORMAT_ISO_8601 = "%Y-%m-%dT%H:%M:%SZ"


def get_sign_string(source, secret):
    secret = bytes(secret.encode('utf8'))
    h = hmac.new(secret, source.encode('utf8'), hashlib.sha1)
    signature = base64.encodebytes(h.digest()).strip()
    return signature


def __pop_standard_urlencode(query):
    ret = parse.urlencode(query)
    ret = ret.replace('+', '%20')
    ret = ret.replace('*', '%2A')
    ret = ret.replace('%7E', '~')
    return ret


def __compose_string_to_sign(method, queries):
    sorted_parameters = sorted(queries.items(), key=lambda queries: queries[0])
    string_to_sign = method + "&%2F&" + parse.quote(__pop_standard_urlencode(sorted_parameters))
    return string_to_sign


def get_sign(paras):
    str_sign = __compose_string_to_sign("GET", paras)
    return get_sign_string(str_sign, access_key_secret + "&")


def my_ecs_action(action_name, **kwargs):
    paras = {
        "SignatureVersion": "1.0",
        "Format": "JSON",
        "Timestamp": time.strftime(FORMAT_ISO_8601, time.gmtime()),
        "AccessKeyId": access_key,
        "SignatureMethod": "HMAC-SHA1",
        "Version": "2014-05-26",
        "Action": action_name,
        "SignatureNonce": str(uuid.uuid4()),
    }
    if kwargs:
        paras.update(kwargs)
    paras['Signature'] = get_sign(paras)
    res = requests.get(
        url=ecs_url,
        params=paras,
    )
    ret = json.loads(res.text)
    print(json.dumps(response, indent=4).decode('unicode_escape'))

if __name__ == "__main__":
    my_ecs_action("DescribeInstanceTypeFamilies", RegionId='cn-beijing', Generation="ecs-1")


