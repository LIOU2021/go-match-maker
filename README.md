# Introduction
这是一个在模拟游戏配对撮合或是交友软体配对的一个project
- demo
    ```go
    go run ./cmd/main.go
    // or
    go run ./_example/main.go
    ```
- 當用戶加入時的roomId屬於當前hub中沒有的，那麼hub會自動增加該room到hub中
# todo
- 新增初始化檢查緩存有無資料的bug
- 測試腳本執行buff比較小的狀態
- config新增
    - 指定房間要取幾個，預設2
    - 每個房間內的成員要取幾個，預設1
- 在server範例中新增ws用作監聽搓合觸發
	- 拿之前寫的chat project來試試看可不可行