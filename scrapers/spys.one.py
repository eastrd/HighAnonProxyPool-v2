from selenium import webdriver
import time


url_template = "http://spys.one/en/free-proxy-list/{NUM}"
url_pages = [""] + [i for i in range(1, 3)]

urls = [url_template.replace("{NUM}", str(page_num)) for page_num in url_pages]

xpath = {
    "ip"        :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[1]/font[2]",
    "port"      :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[1]/font[2]",
    "protocol"  :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[2]/a/font[1]",
    "country"   :   "/html/body/table[2]/tbody/tr[4]/td/table/tbody/tr[*]/td[4]/a",
}

extractor = {
    "ip"        :   lambda text: text.split(":")[0],
    "port"      :   lambda text: text.split(":")[1],
    "protocol"  :   lambda text: text,
    "country"   :   lambda text: text,
}

ips, ports, protocols, countries = [], [], [], []

# Init Chrome Webdriver
driver = webdriver.PhantomJS()

for url in urls:
    driver.get(url)
    ips += [extractor["ip"](ip_element.text)
        for ip_element in driver.find_elements_by_xpath(xpath["ip"])]
    ports += [extractor["port"](port_element.text)
        for port_element in driver.find_elements_by_xpath(xpath["port"])]
    protocols += [extractor["protocol"](protocol_element.text)
        for protocol_element in driver.find_elements_by_xpath(xpath["protocol"])]
    countries += [extractor["country"](country_element.text)
        for country_element in driver.find_elements_by_xpath(xpath["country"])]


print(ips)
print(ports)
print(protocols)
print(countries)

print(len(ips))
print(len(ports))
print(len(protocols))
print(len(countries))
