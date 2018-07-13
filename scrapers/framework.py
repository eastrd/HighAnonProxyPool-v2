from selenium import webdriver
from time import sleep


def configure(url_template, url_pages, xpath, extractor, sleep_time=0):
    '''Function for configuring metrics and scraping the target website via given defined rules.
    Parameters
    ----------
    url_template : string
        Url base string, with {NUM} indicating the placeholder for different pages
    url_pages : list
        A list of string or integer to insert into url_template {NUM}
    xpath : dict
        Containing xpath rules for ip, ports, protocol, country
    extractor : dict
        Used for further extraction for strings extracted from the xpath rules
    sleep_time : integer
        Time in seconds to wait before starting the xpath extraction
    '''

    urls = [url_template.replace("{NUM}", str(page_num)) for page_num in url_pages]
    ips, ports, protocols, countries = [], [], [], []

    # Init Chrome Webdriver
    driver = webdriver.PhantomJS()

    # Set Viewport
    driver.set_window_size(1920, 1080)

    for url in urls:
        driver.get(url)

        sleep(sleep_time)

        ips += [extractor["ip"](ip_element.text)
                for ip_element in driver.find_elements_by_xpath(xpath["ip"])]
        ports += [extractor["port"](port_element.text)
                for port_element in driver.find_elements_by_xpath(xpath["port"])]
        protocols += [extractor["protocol"](protocol_element.text)
                for protocol_element in driver.find_elements_by_xpath(xpath["protocol"])]
        countries += [extractor["country"](country_element.text)
                for country_element in driver.find_elements_by_xpath(xpath["country"])]


    if len(ips) != len(ports) != len(protocols) != len(countries):
        print("Error! Number of data fields collected mismatch: %s %s %s %s" %
              (len(ips), len(ports), len(protocols), len(countries)))
        exit()

    print(ips)
    print(ports)
    print(protocols)
    print(countries)
