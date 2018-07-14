from framework import scrape
from db import save_new_proxy_record
from time import sleep
from loglib import log_print

time_interval = 60 * 10

# Genearlize URL template to hold all proxy urls
url_template = "http://spys.one/en/anonymous-proxy-list/{NUM}"

# url_pages are a list of strings to be inserted into the url template {NUM} location
url_pages = [""] + [i for i in range(1, 3)]


# xpath rules to locate data fields
xpath = {
    "ip"        :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[1]/font[2]",
    "port"      :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[1]/font[2]",
    "protocol"  :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[2]/a/font[1]",
    "country"   :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[4]/a",
}

# Lambda functions to further extract the information
extractor = {
    "ip"        :   lambda text: text.split(":")[0],
    "port"      :   lambda text: text.split(":")[1],
    "protocol"  :   lambda text: text,
    "country"   :   lambda text: text,
}


while True:
    proxy_list_of_dicts = scrape(url_template, url_pages, xpath, extractor)

    for proxy_dict in proxy_list_of_dicts:
        save_new_proxy_record(proxy_dict)
        
    log_print("Finished one round of scraping, sleep for " + str(time_interval) + " seconds")
    sleep(time_interval)
