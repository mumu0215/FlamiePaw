# coding:utf-8
from bs4 import BeautifulSoup
import requests,re
session = "_fofapro_ars_session=51dd72c519918a7db43355a9a24c7596"
header = {
    "Accept":"text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01",
    "Accept-Encoding":"gzip, deflate, br",
    "Accept-Language":"zh-CN,zh;q=0.9",
    "Connection":"keep-alive",
    "User-Agent":"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
    "X-CSRF-Token":"DpraMUR6PuefxdVpDmbZmgW9572Oz4CKSkqLa4u+astRxa+NSW5t0gfjlRB8cESuUrBvrD+zkGA9GFcfEYAVZA==",
    "X-Requested-With":"XMLHttpRequest",
    "Cookie":session,
}
def file_put(str):
    with open("ip.txt","a") as f:
        f.write(str)
 
def spider_ip(url):
    html_doc = requests.get(url = url,headers = header,verify=False).content
    soup = BeautifulSoup(html_doc)

    for link in soup.find_all('a'):
        try:
            if "http" in link.get('href') :
                ip = link.get('href')[2:-2]
                print ip
                result = re.findall(r"\d+\.\d+\.\d+.\d+",ip,re.I)[0]

                print result
                file_put(ip+"\n")
        except Exception:
            pass

 
for i in range(1,1000):
    spider_ip("https://fofa.so/result?full=true&page="+ str(i) +"&qbase64=InQzZ28uY24i")