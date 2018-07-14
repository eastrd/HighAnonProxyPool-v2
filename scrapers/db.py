from pymongo import MongoClient
from loglib import log_print


def save_new_proxy_record(one_proxy_dict_data):
    '''Examine the given single proxy dict object and store into "new" collection
    '''
    log_print("Connect to MongoDB...")
    client = MongoClient("149.28.195.244", 27017, authSource='admin',
                         username="xiaoyu", password="Jianyuhao123$")

    db = client["proxypool"]

    new_storage = db["new"]
    all_storage = db["all"]

    # If historical collection contains the proxy information (Key: ip), then ignore
    if len([i for i in all_storage.find({"ip": one_proxy_dict_data["ip"]})]) > 0:
        print("Found duplicate proxy for %s, ignore..." % one_proxy_dict_data["ip"])
        return

    # Since this ip is new, then store this proxy data in "new" collection, and store only its ip in all storage
    log_print("Store " + one_proxy_dict_data["ip"])
    new_storage.insert(one_proxy_dict_data)
    all_storage.insert({"ip" : one_proxy_dict_data["ip"]})