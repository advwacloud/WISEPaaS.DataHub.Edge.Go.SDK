## Unreleased

### Add
- SendWhenValueChanged 實作

### Fix
- Mac與Linux based 系統斷線重新連線後部分封包遺失

### Change
- FDF 讀取錯誤時不報錯，直接略過
- recover.sqlite 更名為 NODEID_recover.sqlite
- 更改cfgCache.json名稱為 NODEID +_config.json
- 簡化 config.json 結構，並與雲端 config 的狀態配置同步


## 1.0.4
### Fix
- MQTT SSL 正確格式
- recover.sqlite 無資料

### Change
- Config設定檔格式修正為Publish前的封包內容
- 修改設定檔檔名
- 修改取得memory中設定檔中FDF的存取
- 使用 go mod 管理套件