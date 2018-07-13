from scraper import configure


# Genearlize URL template to hold all proxy urls
url_template = "https://hidemy.name/en/proxy-list/?type=hs&anon=4&start={NUM}#list"

# url_pages are a list of strings to be inserted into the url template {NUM} location
url_pages = [64*i for i in range(0, 12)]


# xpath rules to locate data fields
xpath = {
    "ip"        :   '//*[@id="content-section"]/section[1]/div/table/tbody/tr[*]/td[1]',
    "port"      :   '//*[@id = "content-section"]/section[1]/div/table/tbody/tr[*]/td[2]',
    "protocol"  :   '//*[@id="content-section"]/section[1]/div/table/tbody/tr[*]/td[5]',
    "country"   :   '//*[@id="content-section"]/section[1]/div/table/tbody/tr[*]/td[3]/div',
}

# Lambda functions to further extract the information
extractor = {
    "ip"        :   lambda text: text,
    "port"      :   lambda text: text,
    "protocol"  :   lambda text: text,
    "country"   :   lambda text: text.strip(),
}

configure(url_template, url_pages, xpath, extractor, 10)
