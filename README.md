# opencc-api-go
OpenCC API server written in golang

OpenCC API in golang, compatible with HenryQW's TTRSS OpenCC plugin, **with ARM support!**

Listen to 3000 by default.

You can create an instance use `docker run -d -p 3000:3000 --restart=always --name=opencc colinxkl/opencc-api-go`

API call like this:
```bash
curl --request POST \
  --url http://localhost:3000/t2s \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data 'title=繁體中文（中國大陸、澳門、馬新常稱繁體中文，台灣常稱正體中文或繁體中文）
          &content=實際上，兩岸三地的繁體中文出版物並不拘泥於本地標準，有時使用其他字形和異體字是很頻繁的。'
```
Example response:
```
{"title":"繁体中文（中国大陆、澳门、马新常称繁体中文，台湾常称正体中文或繁体中文）","content":"实际上，两岸三地的繁体
中文出版物并不拘泥于本地标准，有时使用其他字形和异体字是很频繁的。"}
```

## Usage
Conversions
There are 10 conversion schemes available in OpenCC:

* s2t: Simplified Chinese to Traditional Chinese 简体到繁体  
* t2s: Traditional Chinese to Simplified Chinese 繁体到简体  
* s2tw: Simplified Chinese to Traditional Chinese (Taiwan Standard) 简体到台湾正体  
* tw2s: Traditional Chinese (Taiwan Standard) to Simplified Chinese 台湾正体到简体  
* s2hk: Simplified Chinese to Traditional Chinese (Hong Kong Standard) 简体到香港繁体（香港小学学习字词表标准）  
* hk2s: Traditional Chinese (Hong Kong Standard) to Simplified Chinese 香港繁体（香港小学学习字词表标准）到简体  
* s2twp: Simplified Chinese to Traditional Chinese (Taiwan Standard) with Taiwanese idiom 简体到繁体（台湾正体标准）并转换爲台湾常用词彙    
* tw2sp: Traditional Chinese (Taiwan Standard) to Simplified Chinese with Mainland Chinese idiom 繁体（台湾正体标准）到简体并转换爲中国大陆常用词彙
* t2tw: Traditional Chinese (OpenCC Standard) to Taiwan Standard 繁体（OpenCC 标准）到台湾正体
* t2hk: Traditional Chinese (OpenCC Standard) to Hong Kong Standard 繁体（OpenCC 标准）到香港繁体（香港小学学习字词表标准）

In order to use t2hk, the address you post to should be http://localhost:3000/t2hk.

Conversion scheme t2s will be used if don't specify any.

Params
title, text to convert, optional
content, text to convert, optional

## JSON API

A generic JSON API is available at `/api/:convert-mode`.

```bash
curl --request POST \
  --url http://localhost:3000/api/t2s \
  --header 'Content-Type: application/json' \
  --data '{"text": "繁體中文"}'
```

Response:
```json
{
  "converted": "繁体中文"
}
```

The `:convert-mode` can be any of the schemes listed above (e.g., `s2t`, `t2s`, `s2tw`, etc.).

## More Documents
* https://github.com/longbridgeapp/opencc#預設配置文件
* https://github.com/HenryQW/OpenCC.henry.wang
