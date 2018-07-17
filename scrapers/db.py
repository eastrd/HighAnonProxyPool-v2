from pymongo import MongoClient
from loglib import log_print


log_print("Initialize MongoDB Connection...")
client = MongoClient("149.28.195.244", 27017, authSource='admin',
                        username="xiaoyu", password="Jianyuhao123$")


def save_new_proxy_record(one_proxy_dict_data):
    '''Examine the given single proxy dict object and store into "new" collection
    '''
    db = client["proxypool"]

    new_storage = db["new"]
    all_storage = db["all"]

    # If historical collection contains the proxy information (Key: ip), then ignore
    if len([i for i in all_storage.find({"ip": one_proxy_dict_data["ip"]})]) > 0:
        log_print("Found duplicate proxy for " +
                  one_proxy_dict_data["ip"] + ", ignore...")
    else:
        log_print("Store " + one_proxy_dict_data["ip"])
        new_storage.insert(one_proxy_dict_data)
        all_storage.insert({"ip" : one_proxy_dict_data["ip"]})
