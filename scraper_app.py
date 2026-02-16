import requests
# import re
from bs4 import BeautifulSoup as soup


#Gettings Input of what you are going to scrape(search)
product = input("What are you scraping today?: ")

#Targeting specific url, page, and generating page halder
url = f"https://www.kupujemprodajem.com/search.php?action=list&data%5Bpage%5D=1&submit%5Bsearch%5D=Tra%C5%BEi&dummy=name&data%5Bkeywords%5D={product}"
page = requests.get(url).text
doc = soup(page, "html.parser")

#Converting page class to string to be able to split and convert it into intiger for page number
page_text = doc.find(class_="this-page")
pages = int(str(page_text).split("/")[-2].split(">")[-1][:-1])
#print(pages)

#loop that lookups items you want to be parsed, in adition adding some new div for better categorizing
for page in range(pages):
    url = f"https://www.kupujemprodajem.com/search.php?action=list&data%5Bpage%5D={page}&submit%5Bsearch%5D=Tra%C5%BEi&dummy=name&data%5Bkeywords%5D={product}"
    page =requests.get(url).text
    doc = soup(page, "html.parser")
    div = doc.find(class_="adListContainer clearfix")

    product_results = div.find_all(text=re.compile(product))
    for item in product_results:
        print(item)
