# SyncClipboard-Go
剪切板同步小助手-PC端

## 说明
这是和[剪切板同步小助手-安卓端](https://github.com/Dawnnnnnn/SyncClipboard-demo)配套的工具，用于实现在全PC平台上的剪切板同步

因为是Go写的，所以理论上可以支持所有架构，自行编译所需架构。Release里是Mac M1的编译成品

运行时需要同目录下有`config.toml`文件，具体参数解释可以看安卓端仓库的说明

## Mac下配置开机启动

```bash
touch ~/Library/LaunchAgents/SyncClipboard.plist
vim ~/Library/LaunchAgents/SyncClipboard.plist
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>KeepAlive</key>
        <true/>
        <key>RunAtLoad</key>
        <true/>
        <key>Label</key>
        <string>SyncClipboard</string>
        <key>ProgramArguments</key>
        <array>
            <string>/Users/dawnnnnnn/Downloads/Tools/SyncClipboard/SyncClipboard-darwin-arm64</string>
        </array>
        <key>WorkingDirectory</key>
        <string>/Users/dawnnnnnn/Downloads/Tools/SyncClipboard</string>
        <key>StandardOutPath</key>
        <string>/Users/dawnnnnnn/Downloads/Tools/SyncClipboard/run-out.log</string>
        <key>StandardErrorPath</key>
        <string>/Users/dawnnnnnn/Downloads/Tools/SyncClipboard/run-err.log</string>
    </dict>
</plist>

```


```bash
launchctl load  ~/Library/LaunchAgents/SyncClipboard.plist
launchctl start ~/Library/LaunchAgents/SyncClipboard.plist
```