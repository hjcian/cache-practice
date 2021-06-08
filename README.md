## Description
- 由於練習重點在快取機制，故先使用滿足 interface 的 `fakedb` 實作來滿足服務運行時取得的資料

## Test
`make test`

## See coverage
`make see-coverage`

## To be improved
- `cache.go` 有重複的處理程序，應有進一步抽象化、reuse codes 的可能
- `cache_test.go` 中，等待 async job 完成後產生效果，應該用 retry/timeout 的思維來測，測試結果可能會比直接用 `time.Sleep` 還更可預測
- 在有快取的情況下、且 cache invalid 時 call async job，也可能會有 cache stampede 的問題，此部分也應再進一步避免
  - 目前想到的是實作一個與 `multiplelock` 類似機制的 `multiplecas`，來支援此情境 (需要再熟悉 `atomic` package 的用法)
  - 且因為可允許 dirty data 先回傳，故若有不 lock 的處理方式，應還可再進一步提升系統吞吐量
- 加入真實的 DB 實作，並使用 `sqlmock` 來針對一般的情境添加測試

## Reference
- [Multiple Lock Based on Input](https://medium.com/@kf99916/multiple-lock-based-on-input-in-golang-74931a3c8230)