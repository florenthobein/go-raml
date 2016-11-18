import marshal, strutils, httpclient
import client_itsyouonline

type
  tokenInfo= object
    username: string

  tokenResp = object #respon of call to access_token endpoint
    access_token: string
    token_type: string
    scope: string
    expires_in: int
    info: tokenInfo

proc getAccessToken*(c: Client, clientID: string, clientSecret: string): string =
  # get access token from itsyouonline
  let q = @["grant_type=client_credentials", "client_id=" & clientID, "client_secret=" & clientSecret]
 
  let resp = c.request("https://itsyou.online/v1/oauth/access_token?" & q.join("&"), "POST")
  if resp.status != "200 OK":
    echo "failed to get access token:", resp.body
    quit()
  
  let tr = to[tokenResp](resp.body)
 
  c.setAuthHeader("token " & tr.access_token)
  return tr.access_token

proc createJWTToken*(c: Client, scopes: openArray[string], auds: openArray[string]): string = 
  # build query string
  # TODO : find a better way
  var uri: string = "https://itsyou.online/v1/oauth/jwt"
  
  if scopes.len > 0:
    uri.add("?scope=" & scopes.join(","))
  
  if auds.len > 0:
    if scopes.len == 0:
      uri.add("?")
    else:
      uri.add("&")
    uri.add("aud=" & auds.join(","))

  return c.request(uri, "GET").body