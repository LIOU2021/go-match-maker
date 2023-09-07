# Introduction
这是一个在模拟游戏配对撮合或是交友软体配对的一个project
- demo
    ```go
    go run ./cmd/main.go
    // or
    go run ./_example/main.go
    ```
- 當用戶加入時的roomId屬於當前hub中沒有的，那麼hub會自動增加該room到hub中
# 结构
- 主要就是透过redis 集合（set）来实现
    - 第一层会以HubName命名一个key，其中会存放著所有房间，比如 "a, b, c, d ..."
        - 第二层会有数个key，根据房间ID命名key，value存放的用户的id
- 其中hub.members是一个map，存放著 uid > user资料
    - 预设情况如果架设分布式架构时，也就是多台server调用```myHub.Run()```，此时```hub.members```的资料会有问题
- 监听撮合成交是透过底下接口
    ```go
    myHub.Notification()
    ```
# todo
- 新增初始化檢查緩存有無資料的bug
- 測試腳本執行buff比較小的狀態
- config新增
    - 指定房間要取幾個，預設2
    - 每個房間內的成員要取幾個，預設1
- 在server範例中新增ws用作監聽搓合觸發
	- 拿之前寫的chat project來試試看可不可行