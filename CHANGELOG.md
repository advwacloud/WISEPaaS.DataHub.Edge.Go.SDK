## 1.0.6
### Fix
- Fix customized timestamp parsing error

## 1.0.5
### Fix
- Mac與Linux based 系統斷線重新連線後部分封包遺失
- SendData ts 從永遠Time.Now() 改為傳入的EdgeData時間

### Change
- 暫時停用 FDF 轉換功能

## 1.0.4
### Fix
- MQTT SSL 正確格式
- recover.sqlite 無資料

### Change
- Config設定檔格式修正為Publish前的封包內容
- 修改設定檔檔名
- 修改取得memory中設定檔中FDF的存取
- 使用 go mod 管理套件