import smtp, times, strutils, httpclient, net, os, osproc

proc getLocalIPAddresses(): string =
  result = ""
  let output = execCmdEx("ip addr show 2>/dev/null || ifconfig 2>/dev/null || cat /proc/net/fib_trie 2>/dev/null")
  if output.exitCode == 0:
    for line in output.output.splitLines():
      if line.contains("inet ") and not line.contains("127.0.0.1") and not line.contains("::1"):
        let parts = line.split()
        for part in parts:
          if part.contains(".") and not part.contains(":"):
            let ip = part.split("/")[0]
            if ip != "127.0.0.1":
              result &= "- " & ip & "\n"
              break
  if result == "":
    result = "未找到有效的本地网络 IP 地址"

proc main() =
  let fromEmail = "768305875@qq.com"
  let toEmail = "768305875@qq.com"
  let authCode = "gpfruabgjebubdad"

  let ipInfo = getLocalIPAddresses()
  let subject = "本地 IP 地址信息 " & format(now(), "yyyy-MM-dd HH:mm")
  let body = "本地联网 IP 地址信息：\n\n" & ipInfo & "\n\n发送时间：" & $now() & "\n来自 Nim 程序"

  let msg = createMessage(subject, body, fromEmail, @[toEmail])

  echo "连接 SMTP..."
  var client = newSmtp(useSsl = true)
  client.connect("smtp.qq.com", Port(465))

  echo "登录..."
  client.auth(fromEmail, authCode)

  echo "发送..."
  client.sendMail(fromEmail, @[toEmail], $msg)
  client.close()

  echo "✅ 已发送"

when isMainModule:
  main()