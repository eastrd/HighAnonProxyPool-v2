from pymongo import MongoClient
from loglib import log_print


log_print("Initialize MongoDB Connection...")
client = MongoClient("149.28.195.244", 27017, authSource='admin',
                        username="xiaoyu", password="Jianyuhao123$")


def make_protocol(protocol):
    '''Convert any protocol fields into HTTP or HTTPS, abort SOCKS-only fields.
    
    Parameter:
    ---------
    protocol: String
    
    Return:
    ---------
    - "HTTP"
    - "HTTPS"
    - None otherwise
    '''
    protocol = protocol.lower()
    if "https" in protocol:
        return "https"
    elif "http" in protocol:
        return "http"
    return None


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
        # Since this ip is new, then store this proxy data in "new" collection, and store only its ip in all storage
        protocol = make_protocol(one_proxy_dict_data["protocol"])
        if not protocol:
            # Ignore if it's SOCKS protocol
            return 
        log_print("Store " + one_proxy_dict_data["ip"])
        new_storage.insert(one_proxy_dict_data)
        all_storage.insert({"ip" : one_proxy_dict_data["ip"]})
