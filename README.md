# GetTitle

### introduction：

 一个旨在渗透初期批量处理资产的工具

目前只做了从url到探测资产title的功能，后期会逐步添加nmap或御剑端口扫描结果处理、分类处理扫描出的service、资产初步处理信息呈现。最后目标是从nmap到xray扫描的一键化渗透工具。

御剑扫描结果导出格式如下：

![御剑](pic/yujian.png)

### Usage:
```
Command:
  -m int
        mode choice:
        	1:parse from url list,-uF Needed;
        	2:parse from port scan file,-pF or -xF Needed
  -t int
        thread (default 15)
  -p string
        proxy setting
  -uF string
        url file name
  -pF string
        yujian port scan file
  -xF string
        nmap output xmlFileName


```


### ToDo：

- [x] 文本处理区别不同系统换行符
- [x] 添加访问代理
- [ ] 处理300的重定向状态码，实现重定向跟踪
- [ ] 链接端口扫描结果处理的工具实现整合



