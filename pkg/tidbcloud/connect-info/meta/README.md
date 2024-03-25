1. `endpoint` contains the id and value of Endpoint Types. value is used for display.
2. `os` contains the id and value of Operating System. value is used for display.
3. `ca` contains each file path of the corresponding operating system.
4. `client` is the list of supported tools/frameworks/drivers/ORMs.
5. `variable` is the default value displayed if the actual connection information cannot be retrieved.
6. `connection` the content of each client. `download_ca` means the client needs an explicit CA certificate but on these operating systems no builtin CA bundle exists. Users need to download the ISRG X1 Root Certificated manually. `doc` means there is an official doc page to refer to for detailed steps. This controls the bottom help message. `type` controls the render style on console.

Each file contains the original content of the client's connection information. It will be parsed and transformed to frontend codes at build time.

1. Placeholders should be substituded with proper values. Placeholders include `${username}`, `${password}`, `${host}`, `${port}`, `${database}` and `${ca_path}`. Except `${ca_path}`, others always exist.
2. When display, ${password} might not be created yet.
3. For ${ca_path}, it might be replace with a string, or a placeholder(e.g. Python driver on Windows).
4. For some clients, the values need be escaped first.

https://pingcap.feishu.cn/docx/U0uedTQrCoHFJ0x7PJTcNF4hnuh