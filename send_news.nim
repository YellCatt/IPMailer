import smtp, times, strutils, httpclient, net, os

proc getLocalIPAddresses(): string =
  result = ""
  for iface in getInterfaces():
    for addr in iface.addresses:
      if addr.family == AfInet or addr.family == AfInet6:
        let ipStr = $addr.addr
        if not ipStr.startsWith("127.") and not ipStr.startsWith("::1"):
          result &= "- " & ipStr & "\n"
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