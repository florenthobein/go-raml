import httpclient, strutils, tables

type
  Client* = object
    baseURI*: string
    hc: HttpClient

const defaultBaseURI = "http://localhost:8080"

proc newClient*(baseURI = defaultBaseURI): Client =
  # creates new client
  var c = Client(baseURI: baseURI, hc: newHttpClient())
  c.hc.headers = newHttpHeaders({ "Content-Type": "application/json" })
  return c

proc setAuthHeader*(c: Client, value: string) =
  c.hc.headers = newHttpHeaders({ "Content-Type": "application/json" })
  c.hc.headers.add("Authorization", value)

proc addQueryParams(url: string, queryParams: Table) : string =
  # add query params to the request URL
  result = url
  if len(queryParams) == 0:
    return

  var qp: seq[string] = @[]
  for k,v  in queryParams.pairs():
    qp.add($k & "=" & $v)

  var sep: string = "?"

  if url.find("?") > 0:
    sep = "&"

  result = url & sep & qp.join("&")


proc request*(c: Client, endpoint: string, httpMethod = "GET", body = "", queryParams: Table[string, string] = initTable[string, string]()): httpclient.Response =
  var url: string = endpoint
  if not url.startsWith("http"):
    url = c.baseURI & url

  url = addQueryParams(url, queryParams)
  return c.hc.request(url, httpMethod, body)